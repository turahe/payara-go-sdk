package payara

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	cfg := &Config{
		AppID:      "test-app",
		AppSecret:  "test-secret",
	}
	c := NewClient(cfg)
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.baseURL == "" {
		t.Error("baseURL should be set by default")
	}
	if c.appID != "test-app" || c.appSecret != "test-secret" {
		t.Error("credentials not set")
	}
	if c.httpClient == nil {
		t.Error("httpClient should be set")
	}
	if c.logger == nil {
		t.Error("logger should be NopLogger when nil")
	}
}

func TestClient_WithTimeout(t *testing.T) {
	c := NewClient(&Config{AppID: "a", AppSecret: "b"})
	c2 := c.WithTimeout(5 * time.Second)
	if c2 == c {
		t.Error("expected new client")
	}
	if c2.httpClient.Timeout != 5*time.Second {
		t.Errorf("timeout: got %v", c2.httpClient.Timeout)
	}
}

func TestClient_WithEnvironment(t *testing.T) {
	c := NewClient(&Config{AppID: "a", AppSecret: "b", BaseURL: "http://custom"})
	c2 := c.WithEnvironment(EnvironmentSandbox)
	if c2.baseURL != "https://sandbox.payara.id:9090" {
		t.Errorf("sandbox baseURL: got %s", c2.baseURL)
	}
	c3 := c2.WithEnvironment(EnvironmentProduction)
	if c3.baseURL != "https://openapi.payara.id:7654" {
		t.Errorf("production baseURL: got %s", c3.baseURL)
	}
}

func TestBaseURLForEnvironment(t *testing.T) {
	tests := []struct {
		env    Environment
		expect string
	}{
		{EnvironmentSandbox, "https://sandbox.payara.id:9090"},
		{EnvironmentProduction, "https://openapi.payara.id:7654"},
		{"", "https://openapi.payara.id:7654"},
	}
	for _, tt := range tests {
		got := BaseURLForEnvironment(tt.env)
		if got != tt.expect {
			t.Errorf("env %q: got %s", tt.env, got)
		}
	}
}

func TestClient_Transfer_Balance(t *testing.T) {
	c := NewClient(&Config{AppID: "a", AppSecret: "b"})
	if c.Transfer() == nil {
		t.Error("Transfer() nil")
	}
	if c.Balance() == nil {
		t.Error("Balance() nil")
	}
}

func TestAPIError_Error(t *testing.T) {
	e := &APIError{Code: "ERR", Message: "msg"}
	if e.Error() != "ERR: msg" {
		t.Errorf("Error(): got %s", e.Error())
	}
	e2 := &APIError{Message: "only"}
	if e2.Error() != "only" {
		t.Errorf("Error(): got %s", e2.Error())
	}
}

// Example of using a mock RoundTripper for unit tests (no real HTTP calls).
func ExampleMockRoundTripper() {
	mock := &MockRoundTripper{
		StatusCode: 200,
		Body:       []byte(`{"success":true,"message":"ok","data":{"merchant_id":"M","balance":1000000,"currency":"IDR","last_updated":"2024-01-01T00:00:00Z","status":"ACTIVE"}}`),
	}
	client := NewClient(&Config{
		AppID:      "test",
		AppSecret:  "test",
		BaseURL:    "https://test.payara.id",
		HTTPClient: &http.Client{Transport: mock},
	})
	// Balance().GetBalance(ctx) would need a prior login; use RoundTripFunc to stub both.
	_ = client
}
