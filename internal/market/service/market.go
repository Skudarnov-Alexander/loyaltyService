package service

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type MarketService struct {
	db market.Repository
}

func New(db market.Repository) *MarketService {
	return &MarketService{
		db: db,
	}
}

func (s MarketService) SaveOrder(ctx context.Context, userID, orderID string) error {
	if err := s.db.InsertOrder(ctx, userID, orderID); err != nil {
		return err
	}

	return nil
}

func (s MarketService) CheckOrder(ctx context.Context, userID, orderID string) (bool, error) {
	isExist, err := s.db.CheckOrder(ctx, userID, orderID)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

func (s MarketService) FetchOrders(ctx context.Context, userID string) ([]model.Order, error) {
	orders, err := s.db.SelectOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s MarketService) FetchBalance(ctx context.Context, userID string) (model.Balance, error) {
	balance, err := s.db.SelectBalance(ctx, userID)
	if err != nil {
		return model.Balance{}, err
	}

	return balance, nil
}

func (s MarketService) FetchWithdrawals(ctx context.Context, userID string) ([]model.Withdrawn, error) {
	withdrawns, err := s.db.SelectWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdrawns, nil
}

func (s MarketService) MakeWithdrawal(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error) {
	b, err := s.db.ProcessWithdrawn(ctx, userID, w)
	if err != nil {
		return model.Balance{}, err
	}

	return b, nil
}
