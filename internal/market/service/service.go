package service

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type Service struct {
	db market.Repository
}

func New(db market.Repository) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) SaveOrder(ctx context.Context, userID, orderID string) error {
	if err := s.db.InsertOrder(ctx, userID, orderID); err != nil {
		return err
	}

	return nil
}

func (s *Service) FetchOrders(ctx context.Context, userID string) ([]model.Order, error) {
	orders, err := s.db.SelectOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Service) FetchBalance(ctx context.Context, userID string) (model.Balance, error) {
	balance, err := s.db.SelectBalance(ctx, userID)
	if err != nil {
		return model.Balance{}, err
	}

	return balance, nil
}

func (s *Service) FetchWithdrawals(ctx context.Context, userID string) ([]model.Withdrawn, error) {
	withdrawns, err := s.db.SelectWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdrawns, nil
}

func (s *Service) MakeWithdrawal(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error) {
	b, err := s.db.ProcessWithdrawn(ctx, userID, w)
	if err != nil {
		return model.Balance{}, err
	}

	return b, nil
}


	

/*
	INSERT INTO orders(order_id, uploaded_at, fk_user_id)
VALUES
(7777, '2016-06-22 19:10:25-07', '48ee3d18-cdf0-4298-9a57-3c2943fd8c31');
	*/