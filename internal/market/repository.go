package market

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type Repository interface {
	InsertOrder(ctx context.Context, userID string, order model.Order) error
}