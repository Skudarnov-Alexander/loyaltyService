package market

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type Repository interface {
	InsertOrder(ctx context.Context, userID, orderID string) error
	SelectOrders(ctx context.Context, userID string) ([]model.Order, error)
	SelectBalance(ctx context.Context, userID string) (model.Balance, error)
	ProcessWithdrawn(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error)
	SelectWithdrawals(ctx context.Context, userID string) ([]model.Withdrawn, error)
}