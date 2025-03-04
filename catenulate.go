package spanx

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Catenulate{})
	httpcaddyfile.RegisterHandlerDirective("replace_request_body", parseCatenulate)
}

// Catenulate implements a handler that replaces the request body
type Catenulate struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Catenulate) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.catenulate",
		New: func() caddy.Module { return new(Catenulate) },
	}
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (c Catenulate) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Create a custom ResponseWriter to capture the response
	crw := &captureResponseWriter{ResponseWriter: w}

	rd, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	c.logger.Debug("Preparing to capture response body",
		zap.ByteString("request", rd))

	// Call the next handler
	err = next.ServeHTTP(crw, r)
	if err != nil {
		return err
	}

	// Replace the request body with the captured response body
	r.Body = io.NopCloser(bytes.NewReader(crw.body))
	r.ContentLength = int64(len(crw.body))

	rd, err = httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	c.logger.Debug("Replacing request body with response body",
		zap.ByteString("request", rd),
		zap.ByteString("crw", crw.body))

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

func (c *Catenulate) Provision(ctx caddy.Context) (err error) {
	c.logger = ctx.Logger(c)
	return c.rp.Provision(ctx)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (c *Catenulate) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// parseCatenulate unmarshals tokens from h into a new Catenulate
func parseCatenulate(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var c Catenulate
	err := c.UnmarshalCaddyfile(h.Dispenser)
	return c, err
}

// Interface guards
var (
	_ caddy.Module                = (*Catenulate)(nil)
	_ caddyhttp.MiddlewareHandler = (*Catenulate)(nil)
	_ caddyfile.Unmarshaler       = (*Catenulate)(nil)
)
