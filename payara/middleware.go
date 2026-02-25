package payara

import (
	"io"
	"net/http"
	"time"
)

// LoggingMiddleware returns a Middleware that logs request and response using the injected Logger.
// If Logger is nil or NopLogger, it no-ops. Use WithMiddleware(LoggingMiddleware(logger)) to add.
func LoggingMiddleware(logger Logger) Middleware {
	if logger == nil {
		return func(next http.RoundTripper) http.RoundTripper {
			return next
		}
	}
	return func(next http.RoundTripper) http.RoundTripper {
		return &loggingRoundTripper{next: next, logger: logger}
	}
}

type loggingRoundTripper struct {
	next   http.RoundTripper
	logger Logger
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	l.logger.Debug("payara request", "method", req.Method, "url", req.URL.String())
	resp, err := l.next.RoundTrip(req)
	if err != nil {
		l.logger.Error("payara request failed", "error", err, "url", req.URL.String())
		return nil, err
	}
	l.logger.Debug("payara response", "status", resp.StatusCode, "url", req.URL.String(), "duration_ms", time.Since(start).Milliseconds())
	return resp, nil
}

// OpenTelemetryMiddleware returns a Middleware that runs the given hook for each request/response.
// Hook can record span, attributes, etc. If hook is nil, the middleware no-ops.
// Example: otelHook could start a span, set attributes from req, then end span with resp/error.
func OpenTelemetryMiddleware(otelHook func(req *http.Request, resp *http.Response, err error)) Middleware {
	if otelHook == nil {
		return func(next http.RoundTripper) http.RoundTripper { return next }
	}
	return func(next http.RoundTripper) http.RoundTripper {
		return &otelRoundTripper{next: next, hook: otelHook}
	}
}

type otelRoundTripper struct {
	next http.RoundTripper
	hook func(*http.Request, *http.Response, error)
}

func (o *otelRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := o.next.RoundTrip(req)
	o.hook(req, resp, err)
	return resp, err
}

// Ensure RoundTripper is implemented
var _ http.RoundTripper = (*loggingRoundTripper)(nil)
var _ http.RoundTripper = (*otelRoundTripper)(nil)

// noopReadCloser is used when we need to replace response body for retry (body already consumed).
type noopReadCloser struct{ io.Reader }

func (n noopReadCloser) Close() error { return nil }
