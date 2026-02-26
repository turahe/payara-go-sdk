//go:build integration

package payara

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/turahe/payara-go-sdk/payara/types"
)

// Integration tests require PAYARA_APP_ID and PAYARA_APP_SECRET (sandbox).
// Run with: go test -tags=integration -v ./payara/...
func TestIntegration_Balance(t *testing.T) {
	appID := os.Getenv("PAYARA_APP_ID")
	appSecret := os.Getenv("PAYARA_APP_SECRET")
	if appID == "" || appSecret == "" {
		t.Skip("set PAYARA_APP_ID and PAYARA_APP_SECRET for integration tests")
	}
	client := NewClient(&Config{
		BaseURL:   BaseURLForEnvironment(EnvironmentSandbox),
		AppID:     appID,
		AppSecret: appSecret,
	}).WithEnvironment(EnvironmentSandbox)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	bal, err := client.Balance().GetBalance(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if bal.Data == nil {
		t.Fatal("balance data nil")
	}
	t.Logf("balance: %d %s", bal.Data.Balance, bal.Data.Currency)
}

func TestIntegration_Disbursement_Status(t *testing.T) {
	appID := os.Getenv("PAYARA_APP_ID")
	appSecret := os.Getenv("PAYARA_APP_SECRET")
	if appID == "" || appSecret == "" {
		t.Skip("set PAYARA_APP_ID and PAYARA_APP_SECRET for integration tests")
	}
	client := NewClient(&Config{
		BaseURL:   BaseURLForEnvironment(EnvironmentSandbox),
		AppID:     appID,
		AppSecret: appSecret,
	}).WithEnvironment(EnvironmentSandbox)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req := types.CreateDisbursementRequest{
		ReferenceID:   "IT-" + time.Now().Format("20060102150405"),
		Amount:        10000,
		BankCode:      "5",
		AccountNumber: "12330922231",
		AccountName:   "Asep",
		Description:   "Integration test",
	}
	createResp, err := client.Transfer().CreateDisbursement(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("created: %s", createResp.Data.TransactionID)
	statusResp, err := client.Transfer().GetDisbursementStatus(ctx, createResp.Data.TransactionID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("status: %s", statusResp.Data.Status)
}
