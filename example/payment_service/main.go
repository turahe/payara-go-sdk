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
	log.Printf("balance: %d %s", bal.Data.Balance, bal.Data.Currency)

	// Create disbursement (use sandbox dummy data from doc: bank_code 5, account 12330922231, name Asep)
	req := types.CreateDisbursementRequest{
		ReferenceID:   "REF-" + time.Now().Format("20060102150405"),
		Amount:        100000, // IDR 100,000 (whole units; do not use float)
		BankCode:      "5",
		AccountNumber: "12330922231",
		AccountName:   "Asep",
		Description:   "Salary payment",
	}
	transferSvc := client.Transfer()
	createResp, err := transferSvc.CreateDisbursement(ctx, req)
	if err != nil {
		log.Fatalf("disbursement: %v", err)
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
