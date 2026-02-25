// Package callback provides an example HTTP handler for Payara disbursement callbacks.
// Doc: https://doc.payara.id/docs/1.0/callback
// Configure your callback URL in the Payara merchant dashboard (Integrations).
package callback

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/turahe/payara-go-sdk/payara/types"
)

// PayaraCallbackHandler is an example HTTP handler for Payara disbursement callbacks.
// Doc: POST with JSON body (transaction_id, amount, status, reference_id, admin_fee, is_refund).
// Required fields: transaction_id, reference_id, status.
// Returns 200 with {"status":"received"} on success; 400 on invalid payload; 500 to trigger Payara retry.
// TODO: Signature validation is not documented by Payara; add verification when documented.
func PayaraCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusBadRequest)
		return
	}

	var payload types.CallbackPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("payara callback: decode error: %v", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields per doc
	if payload.TransactionID == "" || payload.ReferenceID == "" || payload.Status == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing required fields: transaction_id, reference_id, status"})
		return
	}

	// TODO: Validate signature if Payara provides signature verification in future docs

	log.Printf("payara callback: ref=%s txn=%s status=%s amount=%s admin_fee=%s is_refund=%v",
		payload.ReferenceID, payload.TransactionID, payload.Status, payload.Amount, payload.AdminFee, payload.IsRefund)

	// Idempotent processing: update your DB by reference_id, skip if already processed
	// switch payload.Status {
	// case types.CallbackStatusSuccess:
	// 	handleSuccess(payload)
	// case types.CallbackStatusFailed:
	// 	handleFailed(payload)
	// case types.CallbackStatusProcess:
	// 	// no-op or log
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}
