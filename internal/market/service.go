package market

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

//go:generate mockgen -source=service.go -destination=service/mock/mock.go -package=mock

type MarketService interface {
	CheckOrder(ctx context.Context, userID, orderID string) (bool, error)
	SaveOrder(ctx context.Context, userID, orderID string) error
	FetchOrders(ctx context.Context, userID string) ([]model.Order, error)
	FetchBalance(ctx context.Context, userID string) (model.Balance, error)
	FetchWithdrawals(ctx context.Context, userID string) ([]model.Withdrawn, error)
	MakeWithdrawal(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error)
}

type AccrualService interface {
	FetchAccrualStatus(ctx context.Context) error
}

