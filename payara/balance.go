package payara

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/turahe/payara-go-sdk/payara/types"
)

const balancePath = "/api/v1/balance"

// GetBalance sends GET /api/v1/balance. Doc: Get Balance
func (s *balanceService) GetBalance(ctx context.Context) (*types.BalanceResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, s.client.baseURL+balancePath, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := readAll(resp.Body)
	var out types.BalanceResponse
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
