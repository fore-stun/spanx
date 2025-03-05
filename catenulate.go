package spanx

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Catenulate{})
	httpcaddyfile.RegisterHandlerDirective("chain_reverse_proxy", parseCatenulate)
}

// Catenulate implements a handler that replaces the request body
type Catenulate struct {
	logger *zap.Logger
	rp     reverseproxy.Handler
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
	var buffer bytes.Buffer

	// Create a custom ResponseWriter to capture the response
	crw := &captureResponseWriter{
		// logger:         c.logger,
		ResponseWriter: w,
		body:           &buffer,
	}

	rd, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	c.logger.Debug("Preparing to run reverse proxy and capture response body",
		zap.ByteString("request", rd))

	// Invoke the reverse proxy
	err = c.rp.ServeHTTP(crw, r, nil)
	if err != nil {
		return err
	}

	body := buffer.Bytes()
	c.logger.Debug("Extracted body",
		zap.ByteString("body", body))

	// Replace the request body with the captured response body
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))

	// Update headers if necessary
	r.Header.Set("Content-Length", strconv.Itoa(len(body)))
	r.Header.Set("Content-Type", crw.contentType)

	rd, err = httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	c.logger.Debug("Replacing request body with response body",
		zap.ByteString("request", rd),
		zap.ByteString("body", body))

	// Call the next handler
	return next.ServeHTTP(w, r)
}

// captureResponseWriter is a custom ResponseWriter that captures the response body
type captureResponseWriter struct {
	// logger *zap.Logger
	http.ResponseWriter
	body        *bytes.Buffer
	contentType string
}

func (crw *captureResponseWriter) Write(b []byte) (int, error) {
	// crw.logger.Debug("Writing and storing response",
	// 	zap.ByteString("body", b))
	crw.body.Write(b)
	return crw.ResponseWriter.Write(b)
}

func (crw *captureResponseWriter) WriteHeader(statusCode int) {
	crw.contentType = crw.Header().Get("Content-Type")
	crw.ResponseWriter.WriteHeader(statusCode)
}

func (c *Catenulate) Provision(ctx caddy.Context) (err error) {
	c.logger = ctx.Logger(c)
	return c.rp.Provision(ctx)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (c *Catenulate) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	// c.logger.Debug("Creating reverse proxy")
	c.rp = reverseproxy.Handler{}
	err := c.rp.UnmarshalCaddyfile(d)
	if err != nil {
		return err
	}
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
