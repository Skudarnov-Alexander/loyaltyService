package service

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
)

type Service struct {
	db market.Repository
}

func New(db market.Repository) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) SaveOrder(ctx context.Context, orderID string) error {
	return nil
}

func (s *Service) FetchOrders(ctx context.Context) error {
	return nil
}

func (s *Service) FetchBalance(ctx context.Context) error {
	return nil
}

func (s *Service) FetchWithdrawals(ctx context.Context) error {
	return nil
}

func (s *Service) MakeWithdrawal(ctx context.Context) error {
	return nil
}


	

/*
	INSERT INTO orders(order_id, uploaded_at, fk_user_id)
VALUES
(7777, '2016-06-22 19:10:25-07', '48ee3d18-cdf0-4298-9a57-3c2943fd8c31');
	*/