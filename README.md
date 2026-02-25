# Payara API Go SDK

Production-grade Go SDK for [Payara API v1.0](https://doc.payara.id/docs/1.0/). Stateless, concurrency-safe, and microservice-ready.

## Installation

```bash
go get github.com/payara-id/go-sdk/payara
go get github.com/payara-id/go-sdk/payara/types
```

For local development (clone this repo):

```bash
# In your application's go.mod
replace github.com/payara-id/go-sdk => /path/to/payara
```

## Configuration

```go
import (
    "github.com/payara-id/go-sdk/payara"
    "net/http"
    "time"
)

cfg := &payara.Config{
    BaseURL:    payara.BaseURLForEnvironment(payara.EnvironmentSandbox),
    AppID:      os.Getenv("PAYARA_APP_ID"),
    AppSecret:  os.Getenv("PAYARA_APP_SECRET"),
    HTTPClient: &http.Client{Timeout: 30 * time.Second},
    Logger:     myLogger, // inject your logger; nil uses NopLogger
}
client := payara.NewClient(cfg)
```

## Environment setup

| Environment | Base URL |
|-------------|----------|
| Sandbox     | `https://sandbox.payara.id:9090` |
| Production  | `https://openapi.payara.id:7654` |

```go
client = client.WithEnvironment(payara.EnvironmentSandbox)
// or
client = payara.NewClient(&payara.Config{
    BaseURL: payara.BaseURLForEnvironment(payara.EnvironmentProduction),
    AppID:   "...",
    AppSecret: "...",
})
```

Obtain **APP ID** and **APP Secret** from [Payara Merchant Dashboard](https://merchant.payara.id/) → Integrations.

## Login flow and token refresh

- The client **does not** require you to call login manually. The first authenticated request triggers login (POST `/api/v1/login` with `username=app_id`, `password=app_secret`).
- The access token is cached and **refreshed automatically** before expiry (token lifetime 3600s; refresh is triggered 5 minutes before expiry).
- On **401 Unauthorized**, the client retries once after re-login.
- All of this is **thread-safe** (mutex-protected token refresh).

## Required headers

The SDK adds these headers for you on every authenticated request:

- `Authorization: Bearer <access_token>`
- `Content-Type: application/json`
- `Accept: application/json`

Optional: `X-API-Version: 1.0` (documented but not set by default).

## Example usage in a microservice

```go
client := payara.NewClient(cfg).
    WithEnvironment(payara.EnvironmentSandbox).
    WithTimeout(15 * time.Second).
    WithRetryPolicy(payara.DefaultRetryPolicy())

// Balance
bal, err := client.Balance().GetBalance(ctx)
if err != nil {
    return err
}

// Disbursement
req := types.CreateDisbursementRequest{
    ReferenceID:   "REF-UNIQUE-001",
    Amount:         100000, // IDR whole units (min 10_000, max 50_000_000)
    BankCode:       "5",
    AccountNumber:  "12330922231",
    AccountName:    "Asep",
    Description:   "Salary payment",
}
resp, err := client.Transfer().CreateDisbursement(ctx, req)

// Status
status, err := client.Transfer().GetDisbursementStatus(ctx, resp.Data.TransactionID)
```

## Retry strategy

- Use `client.WithRetryPolicy(payara.DefaultRetryPolicy())` to enable retries.
- **Retries** only on **5xx** and **network errors** (exponential backoff).
- Default: max 3 retries, initial backoff 1s, max backoff 30s, multiplier 2.
- Customize with `payara.RetryPolicy{ MaxRetries: 5, Initial: 2*time.Second, ... }`.

## Callback handler

Configure your callback URL in the Payara dashboard. Payara sends a POST with JSON body. Example handler (see `example/callback`):

```go
http.HandleFunc("/callback/payara", callback.PayaraCallbackHandler)
```

Required fields in callback payload: `transaction_id`, `reference_id`, `status`. Return **200** with `{"status":"received"}` for success; **500** triggers Payara retry. **Signature verification** is not documented by Payara; add when/if documented.

## Error handling

```go
resp, err := client.Transfer().CreateDisbursement(ctx, req)
if err != nil {
    var apiErr *payara.APIError
    if errors.As(err, &apiErr) {
        // apiErr.Code, apiErr.Message, apiErr.HTTPStatus, apiErr.RawBody
    }
    return err
}
```

## Money handling

- **Do not use `float64`** for amounts.
- All amounts are **int64** in **IDR whole units** (no decimal; IDR has no minor unit).
- Min disbursement: 10,000 IDR; max: 50,000,000 IDR.

## Production notes

1. Use **Production** base URL and credentials for live traffic.
2. Inject a **logger** (e.g. zerolog, zap) via `Config.Logger` for observability.
3. Use **WithRetryPolicy** for resilience to 5xx and transient network errors.
4. Set **timeouts** with `WithTimeout` or `Config.HTTPClient.Timeout`.
5. Implement **idempotency** for callbacks (key by `reference_id` / `transaction_id`).
6. **ListDisbursement** is not implemented; Payara 1.0 docs do not document a list endpoint. Use **GetDisbursementStatus** by `transaction_id` or `reference_id` instead.

## Package layout

- `payara` – client, config, auth, middleware, retry, errors, transfer, balance.
- `payara/types` – request/response types and enums.
- `example/payment_service` – full payment flow example.
- `example/withdrawal_service` – withdrawal example.
- `example/callback` – example callback HTTP handler.

## License

Use according to your organization’s policy. Payara API terms apply to API usage.
