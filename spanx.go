package spanx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

type JSONFromMultipartForm struct {
}

func (c JSONFromMultipartForm) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
	next caddyhttp.Handler,
) error {
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "multipart/form-data" {
		jsonPayload, err := ConvertFormDataToJSON(r)
		if err != nil {
			http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
			return nil
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(jsonPayload))
		r.ContentLength = int64(len(jsonPayload))
		r.Header.Set("Content-Type", "application/json")
	}

	return next.ServeHTTP(w, r)
}

func ConvertFormDataToJSON(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}

	formValues := make(map[string]interface{})
	for key, values := range r.MultipartForm.Value {
		if len(values) == 1 {
			formValues[key] = values[0]
		} else {
			formValues[key] = values
		}
	}

	jsonPayload, err := json.Marshal(formValues)
	if err != nil {
		return nil, err
	}

	return jsonPayload, nil
}

func init() {
	caddy.RegisterModule(JSONFromMultipartForm{})
}

// CaddyModule returns the Caddy module information.
func (JSONFromMultipartForm) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.spanx",
		New: func() caddy.Module { return new(JSONFromMultipartForm) },
	}
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*JSONFromMultipartForm)(nil)
	_ caddyfile.Unmarshaler       = (*JSONFromMultipartForm)(nil)
)
