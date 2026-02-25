package types

// LoginRequest is the body for POST /api/v1/login.
// Doc: username = app_id, password = app_secret
type LoginRequest struct {
	Username string `json:"username"` // Your app_id
	Password string `json:"password"` // Your app_secret
}

// CreateDisbursementRequest is the body for POST /api/v1/disbursement.
// Amount: min IDR 10,000, max IDR 50,000,000. Use int64 for IDR whole units (no decimal).
// Doc: reference_id must be unique; duplicate rejected.
type CreateDisbursementRequest struct {
	ReferenceID   string `json:"reference_id"`   // Unique transaction reference ID
	Amount        int64  `json:"amount"`        // Disbursement amount in IDR (whole units)
	BankCode      string `json:"bank_code"`      // Recipient bank code
	AccountNumber string `json:"account_number"` // Recipient account number
	AccountName   string `json:"account_name"`   // Recipient account name
	Description   string `json:"description,omitempty"` // Optional transaction description
}

// CheckAccountRequest is the body for POST /api/v1/check-account (account validation).
type CheckAccountRequest struct {
	BankCode      string `json:"bank_code"`      // Bank code to validate
	AccountNumber string `json:"account_number"` // Account number to validate
}

// ListFilter is used for listing disbursements.
// TODO: Payara 1.0 docs do not document a list disbursement endpoint; add fields when API is documented.
type ListFilter struct {
	// ReferenceID optional filter by reference_id
	ReferenceID *string
	// Status optional filter by status
	Status *DisbursementStatus
	// Limit optional page size
	Limit *int
	// Offset optional pagination offset
	Offset *int
}
