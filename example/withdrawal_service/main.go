// Example withdrawal service: check balance and create disbursement.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/turahe/payara-go-sdk/payara"
	"github.com/turahe/payara-go-sdk/payara/types"
)

func main() {
	cfg := &payara.Config{
		BaseURL:   payara.BaseURLForEnvironment(payara.EnvironmentSandbox),
		AppID:     os.Getenv("PAYARA_APP_ID"),
		AppSecret: os.Getenv("PAYARA_APP_SECRET"),
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		log.Fatal("set PAYARA_APP_ID and PAYARA_APP_SECRET")
	}

	client := payara.NewClient(cfg).WithEnvironment(payara.EnvironmentSandbox)
	ctx := context.Background()

	bal, err := client.Balance().GetBalance(ctx)
	if err != nil {
		log.Fatalf("balance: %v", err)
	}
	if bal.Data.Balance < 50000 {
		log.Fatal("insufficient balance")
	}

	req := types.CreateDisbursementRequest{
		ReferenceID:   "WD-" + time.Now().Format("20060102150405"),
		Amount:        50000,
		BankCode:      "5",
		AccountNumber: "12330922231",
		AccountName:   "Asep",
		Description:   "Withdrawal",
	}
	resp, err := client.Transfer().CreateDisbursement(ctx, req)
	if err != nil {
		log.Fatalf("disbursement: %v", err)
	}
	log.Printf("withdrawal submitted: %s", resp.Data.TransactionID)
}
