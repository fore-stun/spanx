package spanx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

type JSONFromMultipartForm struct {
	logger *zap.Logger
}

func (c JSONFromMultipartForm) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
	next caddyhttp.Handler,
) error {
	if r.Method == http.MethodPost && matchesContentType(r, "multipart/form-data") {
		c.logger.Debug("Identified multipart/form-data", zap.String("path", r.URL.Path))
		jsonPayload, err := ConvertFormDataToJSON(r)
		if err != nil {
			http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
			return nil
		}

		r.Body = io.NopCloser(bytes.NewReader(jsonPayload))
		r.ContentLength = int64(len(jsonPayload))
		r.Header.Set("Content-Type", "application/json")
	}

	return next.ServeHTTP(w, r)
}

func matchesContentType(r *http.Request, prefix string) bool {
	// Get the request's Content-Type header
	rct := strings.ToLower(r.Header.Get("Content-Type"))

	// Check if the Content-Type header starts with the given prefix (case-insensitive)
	return strings.HasPrefix(rct, strings.ToLower(prefix))
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
	httpcaddyfile.RegisterHandlerDirective("jaon_from_multipart_form", parseCaddyfile)
}

// CaddyModule returns the Caddy module information.
func (JSONFromMultipartForm) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.spanx",
		New: func() caddy.Module { return new(JSONFromMultipartForm) },
	}
}

func (c *JSONFromMultipartForm) Provision(ctx caddy.Context) (err error) {
	c.logger = ctx.Logger(c)
	return nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (c *JSONFromMultipartForm) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new JSONFromMultipartForm.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var c JSONFromMultipartForm
	err := c.UnmarshalCaddyfile(h.Dispenser)
	return c, err
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*JSONFromMultipartForm)(nil)
	_ caddyfile.Unmarshaler       = (*JSONFromMultipartForm)(nil)
)
