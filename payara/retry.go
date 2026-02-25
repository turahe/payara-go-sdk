package payara

import (
	"net/http"
	"time"
)

// RetryPolicy configures exponential backoff. Retry only on 5xx and network errors.
type RetryPolicy struct {
	MaxRetries int           // Max retry attempts (default 3)
	Initial    time.Duration // Initial backoff (default 1s)
	MaxBackoff time.Duration // Max backoff cap (default 30s)
	Multiplier float64       // Backoff multiplier (default 2)
}

// DefaultRetryPolicy returns a policy suitable for most use cases.
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries: 3,
		Initial:    time.Second,
		MaxBackoff: 30 * time.Second,
		Multiplier: 2,
	}
}

// RetryMiddleware returns a Middleware that retries on 5xx and connection errors with exponential backoff.
func RetryMiddleware(policy *RetryPolicy) Middleware {
	if policy == nil {
		policy = DefaultRetryPolicy()
	}
	maxRetries := policy.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}
	initial := policy.Initial
	if initial <= 0 {
		initial = time.Second
	}
	maxBackoff := policy.MaxBackoff
	if maxBackoff <= 0 {
		maxBackoff = 30 * time.Second
	}
	mult := policy.Multiplier
	if mult <= 0 {
		mult = 2
	}
	return func(next http.RoundTripper) http.RoundTripper {
		return &retryRoundTripper{
			next:       next,
			maxRetries: maxRetries,
			initial:    initial,
			maxBackoff: maxBackoff,
			mult:       mult,
		}
	}
}

type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries int
	initial    time.Duration
	maxBackoff time.Duration
	mult       float64
}

func (r *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response
	backoff := r.initial
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		resp, err := r.next.RoundTrip(req)
		if err != nil {
			lastErr = err
			lastResp = nil
			if !shouldRetryError(err) {
				return nil, err
			}
			if attempt < r.maxRetries {
				sleep(backoff)
				backoff = nextBackoff(backoff, r.maxBackoff, r.mult)
			}
			continue
		}
		if resp.StatusCode < 500 {
			return resp, nil
		}
		lastResp = resp
		lastErr = nil
		if attempt < r.maxRetries {
			resp.Body.Close()
			sleep(backoff)
			backoff = nextBackoff(backoff, r.maxBackoff, r.mult)
		} else {
			return resp, nil
		}
	}
	if lastResp != nil {
		return lastResp, nil
	}
	return nil, lastErr
}

func shouldRetryError(err error) bool {
	// Retry on temporary network errors; could check for net.Error and Temporary()
	return err != nil
}

func nextBackoff(current, max time.Duration, mult float64) time.Duration {
	next := time.Duration(float64(current) * mult)
	if next > max {
		return max
	}
	return next
}

func sleep(d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
	}
}
