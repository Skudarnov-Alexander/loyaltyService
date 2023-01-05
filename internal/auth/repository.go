package auth

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, User *model.User) error
	GetUser(ctx context.Context, username, password string) (*model.User, error)
}
