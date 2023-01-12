package market

import "context"

type MarketService interface {
	SaveOrder(ctx context.Context, orderID string) error
	FetchOrders(ctx context.Context) error
	FetchBalance(ctx context.Context) error
	FetchWithdrawals(ctx context.Context) error
	MakeWithdrawal(ctx context.Context) error
}

type AccrualService interface {
	FetchAccrualStatus(ctx context.Context) error
}

