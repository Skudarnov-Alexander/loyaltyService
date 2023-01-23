package market

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type Repository interface {
	CheckOrder(ctx context.Context, userID, orderID string) (bool, error)
	InsertOrder(ctx context.Context, userID, orderID string) error
	SelectOrders(ctx context.Context, userID string) ([]model.Order, error)
	SelectBalance(ctx context.Context, userID string) (model.Balance, error)
	ProcessWithdrawn(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error)
	SelectWithdrawals(ctx context.Context, userID string) ([]model.Withdrawn, error)
}

type AccrualRepository interface {
	TakeOrdersForProcess(limitOrders int) ([]model.Accrual, error)
	ChangeStatusOrdersForProcess(accruals ...model.Accrual) error
	UpdateStatusProcessedOrders(a model.Accrual) error
	UpdateBalanceProcessedOrders(a model.Accrual) error
}