package payara

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/turahe/payara-go-sdk/payara/types"
)

// wrapWithMiddlewares chains middlewares around the base transport.
func wrapWithMiddlewares(base *http.Client, mws []Middleware) *http.Client {
	if base == nil {
		base = &http.Client{}
	}
	rt := base.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	for i := len(mws) - 1; i >= 0; i-- {
		rt = mws[i](rt)
	}
	clone := *base
	clone.Transport = rt
	return &clone
}

// newJSONRequest builds a POST/GET request with optional JSON body.
func newJSONRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func readAll(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return io.ReadAll(r)
}

// ReadAll reads the body (exported for use by subpackages). Prefer io.ReadAll when in payara package.
func ReadAll(r io.Reader) ([]byte, error) { return readAll(r) }

// ParseErrorResponse tries to decode error body into APIError. Exported for use by subpackages.
func ParseErrorResponse(raw []byte, statusCode int) *APIError {
	return parseErrorResponse(raw, statusCode)
}

// NewJSONRequest builds a request with optional JSON body. Exported for use by subpackages.
func NewJSONRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	return newJSONRequest(ctx, method, url, body)
}

// parseErrorResponse tries to decode error body into ErrorResponse and APIError.
func parseErrorResponse(raw []byte, statusCode int) *APIError {
	var er types.ErrorResponse
	_ = json.Unmarshal(raw, &er)
	return &APIError{
		Code:       er.ErrorCode,
		Message:    er.Message,
		HTTPStatus: statusCode,
		RawBody:    raw,
	}
}
