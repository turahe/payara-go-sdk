package payara

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turahe/payara-go-sdk/payara/types"
)

const (
	disbursementPath = "/api/v1/disbursement"
	checkStatusPath  = "/api/v1/check-status"
)

// CreateDisbursement sends POST /api/v1/disbursement.
// Amount is in IDR whole units (min 10_000, max 50_000_000). reference_id must be unique.
func (s *transferService) CreateDisbursement(ctx context.Context, req types.CreateDisbursementRequest) (*types.CreateDisbursementResponse, error) {
	httpReq, err := newJSONRequest(ctx, http.MethodPost, s.client.baseURL+disbursementPath, req)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := readAll(resp.Body)
	var out types.CreateDisbursementResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, parseErrorResponse(raw, resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseErrorResponse(raw, resp.StatusCode)
	}
	if !out.Success {
		return nil, parseErrorResponse(raw, resp.StatusCode)
	}
	return &out, nil
}

// GetDisbursementStatus sends GET /api/v1/check-status/{id}. Doc: Check Status.
// id can be transaction_id (path) or use GetDisbursementStatusByReference for reference_id (query param).
func (s *transferService) GetDisbursementStatus(ctx context.Context, id string) (*types.DisbursementStatusResponse, error) {
	url := s.client.baseURL + checkStatusPath + "/" + id
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := readAll(resp.Body)
	var out types.DisbursementStatusResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, parseErrorResponse(raw, resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseErrorResponse(raw, resp.StatusCode)
	}
	if !out.Success {
		return nil, parseErrorResponse(raw, resp.StatusCode)
	}
	return &out, nil
}

// ListDisbursement is not implemented. Payara API 1.0 docs do not document a list disbursement endpoint.
// Use GetDisbursementStatus by reference_id or transaction_id instead.
func (s *transferService) ListDisbursement(ctx context.Context, filter types.ListFilter) (*types.DisbursementListResponse, error) {
	return nil, fmt.Errorf("%w", ErrListNotSupported)
}
