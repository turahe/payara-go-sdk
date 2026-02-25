package payara

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

// MockRoundTripper is a simple http.RoundTripper for tests.
// Either set StatusCode/Body/Err for a single fixed response, or set RoundTripFunc for custom behavior.
type MockRoundTripper struct {
	StatusCode    int
	Body          []byte
	Err           error
	RoundTripFunc func(*http.Request) (*http.Response, error)
	mu            sync.Mutex
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	fn := m.RoundTripFunc
	status := m.StatusCode
	body := m.Body
	err := m.Err
	m.mu.Unlock()
	if fn != nil {
		return fn(req)
	}
	if err != nil {
		return nil, err
	}
	if status == 0 {
		status = 200
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}
