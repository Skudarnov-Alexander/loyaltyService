package market

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type MarketService interface {
	SaveOrder(ctx context.Context, userID, orderID string) error
	FetchOrders(ctx context.Context, userID string) ([]model.Order, error)
	FetchBalance(ctx context.Context, userID string) (model.Balance, error)
	FetchWithdrawals(ctx context.Context, userID string) ([]model.Withdrawn, error)
	MakeWithdrawal(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error)
}

type AccrualService interface {
	FetchAccrualStatus(ctx context.Context) error
}

