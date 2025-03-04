package replacebody

import (
	"bytes"
	"io"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(ReplaceBody{})
}

// ReplaceBody implements a handler that replaces the request body
type ReplaceBody struct{}

// CaddyModule returns the Caddy module information.
func (ReplaceBody) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.replace_body",
		New: func() caddy.Module { return new(ReplaceBody) },
	}
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (rb ReplaceBody) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Create a custom ResponseWriter to capture the response
	crw := &captureResponseWriter{ResponseWriter: w}

	// Call the next handler
	err := next.ServeHTTP(crw, r)
	if err != nil {
		return err
	}

	// Replace the request body with the captured response body
	r.Body = io.NopCloser(bytes.NewReader(crw.body))
	r.ContentLength = int64(len(crw.body))

	return nil
}

// captureResponseWriter is a custom ResponseWriter that captures the response body
type captureResponseWriter struct {
	http.ResponseWriter
	body []byte
}

func (crw *captureResponseWriter) Write(b []byte) (int, error) {
	crw.body = append(crw.body, b...)
	return len(b), nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (rb *ReplaceBody) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// Interface guards
var (
	_ caddy.Module                = (*ReplaceBody)(nil)
	_ caddyhttp.MiddlewareHandler = (*ReplaceBody)(nil)
	_ caddyfile.Unmarshaler       = (*ReplaceBody)(nil)
)
