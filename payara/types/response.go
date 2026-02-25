package types

// Meta is common response meta. Doc: timestamp, version (optional)
type Meta struct {
	Timestamp string  `json:"timestamp,omitempty"`
	Version   string  `json:"version,omitempty"`
	RetryAfter *int   `json:"retry_after,omitempty"` // For 429 rate limit
}

// LoginResponseData is the data object from login success response.
// Doc: access_token, token_type "Bearer", expires_in (seconds), merchant_id, merchant_name
type LoginResponseData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // Seconds until expiry
	MerchantID   string `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
}

// GenericAPIResponse is the common envelope: success, message, data, meta.
type GenericAPIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// LoginResponse is the full login response.
type LoginResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *LoginResponseData `json:"data,omitempty"`
	Meta    *Meta              `json:"meta,omitempty"`
}

// ErrorResponse is the documented error format. Doc: success=false, message, error_code
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
	Meta      *Meta  `json:"meta,omitempty"`
}

// CreateDisbursementResponseData is the data object from disbursement creation.
// Doc: transaction_id, reference_id, amount, fee, total_amount, status, bank_code, bank_name,
// account_number, account_name, description, created_at
type CreateDisbursementResponseData struct {
	TransactionID string             `json:"transaction_id"`
	ReferenceID   string             `json:"reference_id"`
	Amount        int64              `json:"amount"`
	Fee           int64              `json:"fee"`
	TotalAmount   int64              `json:"total_amount"`
	Status        DisbursementStatus `json:"status"`
	BankCode      string             `json:"bank_code"`
	BankName      string             `json:"bank_name"`
	AccountNumber string             `json:"account_number"`
	AccountName   string             `json:"account_name"`
	Description   string             `json:"description,omitempty"`
	CreatedAt     string             `json:"created_at"`
}

// CreateDisbursementResponse is the full response for POST /api/v1/disbursement.
type CreateDisbursementResponse struct {
	Success bool                            `json:"success"`
	Message string                          `json:"message"`
	Data    *CreateDisbursementResponseData `json:"data,omitempty"`
	Meta    *Meta                           `json:"meta,omitempty"`
}

// DisbursementStatusData is the data object from check-status.
// Doc: transaction_id, reference_id, status, amount, fee, total_amount, bank_code, bank_name,
// account_number, account_name, description, created_at, processed_at, failure_reason (if FAILED)
type DisbursementStatusData struct {
	TransactionID  string             `json:"transaction_id"`
	ReferenceID    string             `json:"reference_id"`
	Status         DisbursementStatus `json:"status"`
	Amount         int64             `json:"amount"`
	Fee            int64             `json:"fee"`
	TotalAmount    int64             `json:"total_amount"`
	BankCode       string            `json:"bank_code"`
	BankName       string            `json:"bank_name"`
	AccountNumber  string            `json:"account_number"`
	AccountName    string            `json:"account_name"`
	Description    string            `json:"description,omitempty"`
	CreatedAt      string            `json:"created_at"`
	ProcessedAt    string            `json:"processed_at,omitempty"`
	FailureReason  *string           `json:"failure_reason,omitempty"`
}

// DisbursementStatusResponse is the full response for GET /api/v1/check-status
type DisbursementStatusResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    *DisbursementStatusData  `json:"data,omitempty"`
	Meta    *Meta                    `json:"meta,omitempty"`
}

// BalanceData is the data object from GET /api/v1/balance.
// Doc: merchant_id, balance, currency (IDR), last_updated, status (ACTIVE|SUSPENDED|BLOCKED)
type BalanceData struct {
	MerchantID  string        `json:"merchant_id"`
	Balance     int64         `json:"balance"` // IDR whole units
	Currency    string        `json:"currency"`
	LastUpdated string        `json:"last_updated"`
	Status      AccountStatus `json:"status"`
}

// BalanceResponse is the full response for GET /api/v1/balance
type BalanceResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *BalanceData `json:"data,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
}

// DisbursementListResponse is the response for list disbursement.
// TODO: Payara 1.0 docs do not document list endpoint; structure may change when API is available.
type DisbursementListResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    []DisbursementStatusData    `json:"data,omitempty"`
	Meta    *Meta                       `json:"meta,omitempty"`
}

// CallbackPayload is the POST body sent by Payara to the configured callback URL.
// Doc: transaction_id, amount (string IDR), status (Success|Failed|Process), reference_id, admin_fee (string), is_refund (bool)
// No signature verification documented; add if Payara provides it later.
type CallbackPayload struct {
	TransactionID string         `json:"transaction_id"`
	Amount        string         `json:"amount"`   // String in IDR
	Status        CallbackStatus `json:"status"`
	ReferenceID   string         `json:"reference_id"`
	AdminFee      string         `json:"admin_fee"` // String
	IsRefund      bool           `json:"is_refund"` // true = failed refund, false = regular failed
}
