package types

import (
	"encoding/json"
	"testing"
)

func TestFlexString_UnmarshalJSON_number(t *testing.T) {
	var f FlexString
	err := json.Unmarshal([]byte("206"), &f)
	if err != nil {
		t.Fatal(err)
	}
	if f != "206" {
		t.Errorf("got %q, want \"206\"", f)
	}
}

func TestFlexString_UnmarshalJSON_string(t *testing.T) {
	var f FlexString
	err := json.Unmarshal([]byte(`"206"`), &f)
	if err != nil {
		t.Fatal(err)
	}
	if f != "206" {
		t.Errorf("got %q, want \"206\"", f)
	}
}

func TestBalanceAmount_UnmarshalJSON_number(t *testing.T) {
	var b BalanceAmount
	err := json.Unmarshal([]byte("999793000"), &b)
	if err != nil {
		t.Fatal(err)
	}
	if b != 999793000 {
		t.Errorf("got %d, want 999793000", b)
	}
}

func TestBalanceAmount_UnmarshalJSON_stringWithDots(t *testing.T) {
	var b BalanceAmount
	err := json.Unmarshal([]byte(`"999.793.000"`), &b)
	if err != nil {
		t.Fatal(err)
	}
	if b != 999793000 {
		t.Errorf("got %d, want 999793000", b)
	}
}

func TestBalanceAmount_UnmarshalJSON_stringWithCommas(t *testing.T) {
	var b BalanceAmount
	err := json.Unmarshal([]byte(`"1,000,000"`), &b)
	if err != nil {
		t.Fatal(err)
	}
	if b != 1000000 {
		t.Errorf("got %d, want 1000000", b)
	}
}

func TestLoginResponse_Unmarshal_realAPI(t *testing.T) {
	raw := `{"success":true,"message":"Login successful","data":{"access_token":"eyJ...","token_type":"Bearer","expires_in":3599.40778,"merchant_id":206,"merchant_name":"Test Merchant"},"meta":{"timestamp":"2024-01-15T10:30:00Z"}}`
	var out LoginResponse
	err := json.Unmarshal([]byte(raw), &out)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Success || out.Data == nil {
		t.Fatalf("success=%v data=%v", out.Success, out.Data)
	}
	if out.Data.ExpiresIn != 3599.40778 {
		t.Errorf("ExpiresIn = %v, want 3599.40778", out.Data.ExpiresIn)
	}
	if out.Data.MerchantID != "206" {
		t.Errorf("MerchantID = %q, want \"206\"", out.Data.MerchantID)
	}
	if out.Data.MerchantName != "Test Merchant" {
		t.Errorf("MerchantName = %q", out.Data.MerchantName)
	}
}

func TestBalanceResponse_Unmarshal_stringBalance(t *testing.T) {
	raw := `{"success":true,"message":"Balance retrieved successfully","data":{"merchant_id":206,"balance":"999.793.000","currency":"IDR","last_updated":"2026-02-23T22:49:14+07:00","status":"ACTIVE"},"meta":{"timestamp":"2026-02-23T22:49:14+07:00"}}`
	var out BalanceResponse
	err := json.Unmarshal([]byte(raw), &out)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Success || out.Data == nil {
		t.Fatalf("success=%v data=%v", out.Success, out.Data)
	}
	if out.Data.Balance != 999793000 {
		t.Errorf("Balance = %d, want 999793000", out.Data.Balance)
	}
	if out.Data.MerchantID != "206" {
		t.Errorf("MerchantID = %q, want \"206\"", out.Data.MerchantID)
	}
	if out.Data.Currency != "IDR" || out.Data.Status != AccountStatusActive {
		t.Errorf("Currency=%q Status=%q", out.Data.Currency, out.Data.Status)
	}
}

func TestErrorResponse_Unmarshal(t *testing.T) {
	raw := `{"success":false,"message":"Insufficient balance","error_code":"INSUFFICIENT_BALANCE"}`
	var out ErrorResponse
	err := json.Unmarshal([]byte(raw), &out)
	if err != nil {
		t.Fatal(err)
	}
	if out.Success != false || out.ErrorCode != "INSUFFICIENT_BALANCE" || out.Message != "Insufficient balance" {
		t.Errorf("got %+v", out)
	}
}

func TestCallbackPayload_Unmarshal(t *testing.T) {
	raw := `{"transaction_id":"100028355123792503","amount":"10000","status":"Success","reference_id":"REFID153966210","admin_fee":"3500","is_refund":false}`
	var out CallbackPayload
	err := json.Unmarshal([]byte(raw), &out)
	if err != nil {
		t.Fatal(err)
	}
	if out.TransactionID != "100028355123792503" || out.ReferenceID != "REFID153966210" || out.Status != CallbackStatusSuccess {
		t.Errorf("got %+v", out)
	}
	if out.Amount != "10000" || out.AdminFee != "3500" || out.IsRefund != false {
		t.Errorf("Amount=%q AdminFee=%q IsRefund=%v", out.Amount, out.AdminFee, out.IsRefund)
	}
}
