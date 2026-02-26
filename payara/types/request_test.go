package types

import (
	"encoding/json"
	"testing"
)

func TestLoginRequest_JSON(t *testing.T) {
	req := LoginRequest{Username: "app-id", Password: "secret"}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var out LoginRequest
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Username != req.Username || out.Password != req.Password {
		t.Errorf("round-trip: got Username=%q Password=%q", out.Username, out.Password)
	}
}

func TestCreateDisbursementRequest_JSON(t *testing.T) {
	req := CreateDisbursementRequest{
		ReferenceID:   "REF-001",
		Amount:        100000,
		BankCode:      "5",
		AccountNumber: "1234567890",
		AccountName:   "John",
		Description:   "Payment",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var out CreateDisbursementRequest
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.ReferenceID != req.ReferenceID || out.Amount != req.Amount || out.BankCode != req.BankCode {
		t.Errorf("round-trip: got %+v", out)
	}
}

func TestCheckAccountRequest_JSON(t *testing.T) {
	req := CheckAccountRequest{BankCode: "5", AccountNumber: "1234567890"}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var out CheckAccountRequest
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.BankCode != req.BankCode || out.AccountNumber != req.AccountNumber {
		t.Errorf("round-trip: got %+v", out)
	}
}
