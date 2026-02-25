// Package types defines request/response and enum types for Payara API v1.0.
// All types follow the official documentation at https://doc.payara.id/docs/1.0/
package types

// DisbursementStatus represents transaction status from API (check-status, disbursement response).
// Doc: PROCESS | SUCCESS | FAILED
type DisbursementStatus string

const (
	DisbursementStatusProcess DisbursementStatus = "PROCESS"
	DisbursementStatusSuccess DisbursementStatus = "SUCCESS"
	DisbursementStatusFailed  DisbursementStatus = "FAILED"
)

// CallbackStatus represents status in callback payload. Doc uses "Success", "Failed", "Process".
type CallbackStatus string

const (
	CallbackStatusSuccess CallbackStatus = "Success"
	CallbackStatusFailed  CallbackStatus = "Failed"
	CallbackStatusProcess CallbackStatus = "Process"
)

// AccountStatus represents merchant/account status. Doc: ACTIVE | SUSPENDED | BLOCKED
type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "ACTIVE"
	AccountStatusSuspended AccountStatus = "SUSPENDED"
	AccountStatusBlocked   AccountStatus = "BLOCKED"
)

// Environment is Sandbox or Production. Doc: Environmental Information
type Environment string

const (
	EnvironmentSandbox    Environment = "sandbox"
	EnvironmentProduction Environment = "production"
)
