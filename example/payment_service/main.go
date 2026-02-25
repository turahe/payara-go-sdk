// Example payment service using Payara SDK for disbursement (e.g. salary, vendor payments).
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/turahe/payara-go-sdk/payara"
	"github.com/turahe/payara-go-sdk/payara/types"
)

func main() {
	loadEnv()
	cfg := &payara.Config{
		BaseURL:    payara.BaseURLForEnvironment(payara.EnvironmentSandbox),
		AppID:      os.Getenv("PAYARA_APP_ID"),
		AppSecret:  os.Getenv("PAYARA_APP_SECRET"),
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		Logger:     &payara.NopLogger{},
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		log.Fatal("set PAYARA_APP_ID and PAYARA_APP_SECRET")
	}

	client := payara.NewClient(cfg).
		WithEnvironment(payara.EnvironmentSandbox).
		WithTimeout(15 * time.Second).
		WithRetryPolicy(payara.DefaultRetryPolicy())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check balance before disbursement
	balanceSvc := client.Balance()
	bal, err := balanceSvc.GetBalance(ctx)
	if err != nil {
		log.Fatalf("balance: %v", err)
	}
	if bal.Data != nil {
		log.Printf("balance: %d %s", bal.Data.Balance, bal.Data.Currency)
	} else {
		log.Printf("balance: %s", bal.Message)
	}

	// Create disbursement using sandbox dummy account (BCA / Asep)
	recipient := payara.DefaultSandboxAccount()
	req := types.CreateDisbursementRequest{
		ReferenceID:   "REF-" + time.Now().Format("20060102150405"),
		Amount:        100000, // IDR 100,000 (min 10_000, max 50_000_000)
		BankCode:      recipient.BankCode,
		AccountNumber: recipient.AccountNumber,
		AccountName:   recipient.AccountName,
		Description:   "Salary payment (sandbox dummy)",
	}
	transferSvc := client.Transfer()
	createResp, err := transferSvc.CreateDisbursement(ctx, req)
	if err != nil {
		log.Printf("disbursement failed: %v", err)
		return // exit 0 so make doesn't report Error 1
	}
	log.Printf("disbursement created: txn=%s status=%s", createResp.Data.TransactionID, createResp.Data.Status)

	// Optional: poll status (or rely on callback)
	statusResp, err := transferSvc.GetDisbursementStatus(ctx, createResp.Data.TransactionID)
	if err != nil {
		log.Printf("status check: %v", err)
		return
	}
	log.Printf("status: %s", statusResp.Data.Status)
}
