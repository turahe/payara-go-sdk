package types

import "testing"

func TestDisbursementStatus_constants(t *testing.T) {
	if DisbursementStatusProcess != "PROCESS" {
		t.Errorf("DisbursementStatusProcess = %q, want PROCESS", DisbursementStatusProcess)
	}
	if DisbursementStatusSuccess != "SUCCESS" {
		t.Errorf("DisbursementStatusSuccess = %q, want SUCCESS", DisbursementStatusSuccess)
	}
	if DisbursementStatusFailed != "FAILED" {
		t.Errorf("DisbursementStatusFailed = %q, want FAILED", DisbursementStatusFailed)
	}
}

func TestCallbackStatus_constants(t *testing.T) {
	if CallbackStatusSuccess != "Success" {
		t.Errorf("CallbackStatusSuccess = %q, want Success", CallbackStatusSuccess)
	}
	if CallbackStatusFailed != "Failed" {
		t.Errorf("CallbackStatusFailed = %q, want Failed", CallbackStatusFailed)
	}
	if CallbackStatusProcess != "Process" {
		t.Errorf("CallbackStatusProcess = %q, want Process", CallbackStatusProcess)
	}
}

func TestAccountStatus_constants(t *testing.T) {
	if AccountStatusActive != "ACTIVE" {
		t.Errorf("AccountStatusActive = %q, want ACTIVE", AccountStatusActive)
	}
	if AccountStatusSuspended != "SUSPENDED" {
		t.Errorf("AccountStatusSuspended = %q, want SUSPENDED", AccountStatusSuspended)
	}
	if AccountStatusBlocked != "BLOCKED" {
		t.Errorf("AccountStatusBlocked = %q, want BLOCKED", AccountStatusBlocked)
	}
}

func TestEnvironment_constants(t *testing.T) {
	if EnvironmentSandbox != "sandbox" {
		t.Errorf("EnvironmentSandbox = %q, want sandbox", EnvironmentSandbox)
	}
	if EnvironmentProduction != "production" {
		t.Errorf("EnvironmentProduction = %q, want production", EnvironmentProduction)
	}
}
