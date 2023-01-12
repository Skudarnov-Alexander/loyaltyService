package auth

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, u model.User) error
	GetUser(ctx context.Context, username string) (model.User, error)
}
