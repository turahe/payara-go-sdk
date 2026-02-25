package payara

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/turahe/payara-go-sdk/payara/types"
)

func TestTransferService_CreateDisbursement_Mock(t *testing.T) {
	// 1) Mock login response
	loginBody := []byte(`{"success":true,"message":"ok","data":{"access_token":"tok","token_type":"Bearer","expires_in":3600,"merchant_id":"M1","merchant_name":"Test"}}`)
	// 2) Mock disbursement response
	disbBody := []byte(`{"success":true,"message":"ok","data":{"transaction_id":"T1","reference_id":"R1","amount":100000,"fee":2500,"total_amount":102500,"status":"PROCESS","bank_code":"5","bank_name":"BCA","account_number":"123","account_name":"A","description":"","created_at":"2024-01-01T00:00:00Z"}}`)

	var callCount int
	mock := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			var body []byte
			if req.URL.Path == "/api/v1/login" {
				body = loginBody
			} else if req.URL.Path == "/api/v1/disbursement" {
				body = disbBody
			} else {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		},
	}
	client := NewClient(&Config{
		AppID:      "app",
		AppSecret:  "secret",
		BaseURL:    "https://test.payara.id",
		HTTPClient: &http.Client{Transport: mock},
	})
	ctx := context.Background()
	req := types.CreateDisbursementRequest{
		ReferenceID:   "R1",
		Amount:        100000,
		BankCode:      "5",
		AccountNumber: "123",
		AccountName:   "A",
	}
	resp, err := client.Transfer().CreateDisbursement(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data == nil || resp.Data.TransactionID != "T1" || resp.Data.Status != types.DisbursementStatusProcess {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
	if callCount < 2 {
		t.Errorf("expected at least login + disbursement calls, got %d", callCount)
	}
}
