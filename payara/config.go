package payara

import (
	"net/http"
	"time"
)

// Config holds client configuration. BaseURL can be set directly or via WithEnvironment.
type Config struct {
	BaseURL     string
	AppID       string
	AppSecret   string
	HTTPClient  *http.Client
	Middlewares []Middleware
	Logger      Logger
	// RetryPolicy if set is applied when using WithRetryPolicy; optional at construction
	RetryPolicy *RetryPolicy
}

// withDefaults applies default base URL, HTTP client, and middlewares.
func (c *Config) withDefaults() *Config {
	if c.BaseURL == "" {
		c.BaseURL = BaseURLForEnvironment(EnvironmentProduction)
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	if c.Middlewares == nil {
		c.Middlewares = []Middleware{}
	}
	if c.Logger == nil {
		c.Logger = &NopLogger{}
	}
	return c
}
