// Package payara provides a production-grade Go SDK for Payara API v1.0.
// See https://doc.payara.id/docs/1.0/
package payara

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Client is the main API client. It is stateless with respect to request data
// but holds token state for auth. Safe for concurrent use via mutex for token refresh.
type Client struct {
	baseURL     string
	appID       string
	appSecret   string
	accessToken string
	tokenExpiry time.Time
	httpClient  *http.Client
	middlewares []Middleware
	logger      Logger
	auth        *authState
	mu          sync.Mutex
}

// authState holds token and mutex for refresh. Used by Client internally.
type authState struct {
	token   string
	expiry  time.Time
	refresh func(context.Context) error
}

// Logger is the injectable logger interface. Do not hardcode; inject from caller.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// Middleware wraps an HTTP RoundTripper for logging, metrics, etc.
type Middleware func(next http.RoundTripper) http.RoundTripper

// NewClient creates a client with the given config. Uses defaults for nil config.
func NewClient(cfg *Config) *Client {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg = cfg.withDefaults()
	client := &Client{
		baseURL:     cfg.BaseURL,
		appID:       cfg.AppID,
		appSecret:   cfg.AppSecret,
		httpClient:  cfg.HTTPClient,
		middlewares: cfg.Middlewares,
		logger:      cfg.Logger,
	}
	client.auth = &authState{refresh: client.login}
	client.httpClient = wrapWithMiddlewares(client.httpClient, client.middlewares)
	return client
}

// WithTimeout returns a new Config with the given HTTP timeout (copy of config if needed).
// For modifying an existing client, use Config and NewClient.
func (c *Client) WithTimeout(d time.Duration) *Client {
	next := c.httpClient
	if next == nil {
		next = &http.Client{}
	}
	clone := *next
	clone.Timeout = d
	return &Client{
		baseURL:     c.baseURL,
		appID:       c.appID,
		appSecret:   c.appSecret,
		accessToken: c.accessToken,
		tokenExpiry: c.tokenExpiry,
		httpClient:  &clone,
		middlewares: c.middlewares,
		logger:      c.logger,
		auth:        c.auth,
		mu:          c.mu,
	}
}

// WithRetryPolicy returns a new Client with retry middleware prepended.
func (c *Client) WithRetryPolicy(policy *RetryPolicy) *Client {
	middlewares := make([]Middleware, 0, len(c.middlewares)+1)
	middlewares = append(middlewares, RetryMiddleware(policy))
	middlewares = append(middlewares, c.middlewares...)
	hc := wrapWithMiddlewares(c.httpClient, middlewares)
	return &Client{
		baseURL:     c.baseURL,
		appID:       c.appID,
		appSecret:   c.appSecret,
		accessToken: c.accessToken,
		tokenExpiry: c.tokenExpiry,
		httpClient:  hc,
		middlewares: middlewares,
		logger:      c.logger,
		auth:        c.auth,
		mu:          c.mu,
	}
}

// WithMiddleware returns a new Client with the given middleware appended.
func (c *Client) WithMiddleware(m Middleware) *Client {
	middlewares := make([]Middleware, len(c.middlewares), len(c.middlewares)+1)
	copy(middlewares, c.middlewares)
	middlewares = append(middlewares, m)
	hc := wrapWithMiddlewares(c.httpClient, middlewares)
	return &Client{
		baseURL:     c.baseURL,
		appID:       c.appID,
		appSecret:   c.appSecret,
		accessToken: c.accessToken,
		tokenExpiry: c.tokenExpiry,
		httpClient:  hc,
		middlewares: middlewares,
		logger:      c.logger,
		auth:        c.auth,
		mu:          c.mu,
	}
}

// WithEnvironment returns a new Client with base URL set for the given environment.
func (c *Client) WithEnvironment(env Environment) *Client {
	baseURL := BaseURLForEnvironment(env)
	return &Client{
		baseURL:     baseURL,
		appID:       c.appID,
		appSecret:   c.appSecret,
		accessToken: c.accessToken,
		tokenExpiry: c.tokenExpiry,
		httpClient:  c.httpClient,
		middlewares: c.middlewares,
		logger:      c.logger,
		auth:        c.auth,
		mu:          c.mu,
	}
}

// BaseURL returns the configured API base URL.
func (c *Client) BaseURL() string { return c.baseURL }

// DoRequest adds Authorization and performs the request. Refreshes token on 401 and retries once.
func (c *Client) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.doRequest(ctx, req)
}

// Transfer returns the TransferService implementation.
func (c *Client) Transfer() TransferService {
	return &transferService{client: c}
}

// Balance returns the BalanceService implementation.
func (c *Client) Balance() BalanceService {
	return &balanceService{client: c}
}

// Ensure Client implements optional interfaces at compile time.
var (
	_ TransferService = (*transferService)(nil)
	_ BalanceService  = (*balanceService)(nil)
)
