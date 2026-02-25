package payara

import (
	"context"

	"github.com/turahe/payara-go-sdk/payara/types"
)

// TransferService provides disbursement and status operations. Doc: Disbursement, Check Status
type TransferService interface {
	CreateDisbursement(ctx context.Context, req types.CreateDisbursementRequest) (*types.CreateDisbursementResponse, error)
	GetDisbursementStatus(ctx context.Context, id string) (*types.DisbursementStatusResponse, error)
	ListDisbursement(ctx context.Context, filter types.ListFilter) (*types.DisbursementListResponse, error)
}

// BalanceService provides balance inquiry. Doc: Get Balance
type BalanceService interface {
	GetBalance(ctx context.Context) (*types.BalanceResponse, error)
}

// transferService implements TransferService
type transferService struct {
	client *Client
}

// balanceService implements BalanceService
type balanceService struct {
	client *Client
}
