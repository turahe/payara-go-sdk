package payara

import "errors"

// APIError is the structured error for API failures. Doc: success=false, message, error_code
type APIError struct {
	Code       string // error_code from response
	Message    string
	HTTPStatus int
	RawBody    []byte
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return e.Code + ": " + e.Message
	}
	return e.Message
}

// ErrListNotSupported is returned by ListDisbursement. Payara 1.0 docs do not document a list disbursement endpoint.
var ErrListNotSupported = errors.New("payara: list disbursement endpoint not documented in API 1.0")
