package yourmodule

import (
    "bytes"
    "encoding/json"
    "mime/multipart"
    "net/http"

    "github.com/caddyserver/caddy/v2"
    "github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

type ConvertToJSON struct {
    Next    caddyhttp.Handler
}

func (c ConvertToJSON) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
    if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "multipart/form-data" {
        // Parse the multipart form
        err := r.ParseMultipartForm(32 << 20) // Set an appropriate max memory value
        if err != nil {
            // Handle parsing error
            http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
            return nil
        }

        // Convert form values to JSON
        formValues := make(map[string]interface{})
        for key, values := range r.MultipartForm.Value {
            if len(values) == 1 {
                formValues[key] = values[0]
            } else {
                formValues[key] = values
            }
        }

        // Convert form values to JSON payload
        jsonPayload, err := json.Marshal(formValues)
        if err != nil {
            // Handle JSON conversion error
            http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
            return nil
        }

        // Set the request body to the JSON payload
        r.Body = ioutil.NopCloser(bytes.NewReader(jsonPayload))
        r.ContentLength = int64(len(jsonPayload))
        r.Header.Set("Content-Type", "application/json")
    }

    // Call the next handler in the chain
    return c.Next.ServeHTTP(w, r)
}

func init() {
    caddy.RegisterModule(ConvertToJSON{})
}

// CaddyModule returns the Caddy module information.
func (ConvertToJSON) CaddyModule() caddy.ModuleInfo {
    return caddy.ModuleInfo{
        ID:  "http.handlers.yourmodule",
        New: func() caddy.Module { return new(ConvertToJSON) },
    }
}
