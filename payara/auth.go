package payara

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/turahe/payara-go-sdk/payara/types"
)

const (
	loginPath   = "/api/v1/login"
	tokenBuffer = 5 * time.Minute // Refresh token before expiry
)

var errUnauthorized = &APIError{Code: "UNAUTHORIZED", Message: "invalid or expired token"}

// login performs POST /api/v1/login and updates client token. Doc: username=app_id, password=app_secret
func (c *Client) login(ctx context.Context) error {
	body := types.LoginRequest{Username: c.appID, Password: c.appSecret}
	req, err := newJSONRequest(ctx, http.MethodPost, c.baseURL+loginPath, body)
	if err != nil {
		return err
	}
	// Login does not use Bearer; only subsequent API calls do
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, _ := readAll(resp.Body)
	var loginResp types.LoginResponse
	if err := json.Unmarshal(raw, &loginResp); err != nil {
		return &APIError{Message: "login response decode failed", RawBody: raw, HTTPStatus: resp.StatusCode}
	}
	if !loginResp.Success || loginResp.Data == nil {
		code := ""
		var errResp types.ErrorResponse
		_ = json.Unmarshal(raw, &errResp)
		if errResp.ErrorCode != "" {
			code = errResp.ErrorCode
		}
		return &APIError{Code: code, Message: loginResp.Message, HTTPStatus: resp.StatusCode, RawBody: raw}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{Message: loginResp.Message, HTTPStatus: resp.StatusCode, RawBody: raw}
	}

	c.accessToken = loginResp.Data.AccessToken
	expirySec := loginResp.Data.ExpiresIn
	if expirySec <= 0 {
		expirySec = 3600
	}
	c.tokenExpiry = time.Now().Add(time.Duration(expirySec) * time.Second)
	if c.auth != nil {
		c.auth.token = c.accessToken
		c.auth.expiry = c.tokenExpiry
	}
	return nil
}

// getAuthHeader returns the Authorization header value. Call while holding c.mu or after ensureToken.
func (c *Client) getAuthHeader() string {
	return "Bearer " + strings.TrimSpace(c.accessToken)
}

// ensureToken refreshes token if expired or missing. Safe for concurrent use.
func (c *Client) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	needRefresh := c.accessToken == "" || time.Now().Add(tokenBuffer).After(c.tokenExpiry)
	c.mu.Unlock()
	if !needRefresh {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.accessToken != "" && time.Now().Add(tokenBuffer).Before(c.tokenExpiry) {
		return nil
	}
	return c.login(ctx)
}

// doRequest adds auth and performs the request. Refreshes token on 401 and retries once.
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		if err := c.login(ctx); err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.getAuthHeader())
		return c.httpClient.Do(req)
	}
	return resp, nil
}
