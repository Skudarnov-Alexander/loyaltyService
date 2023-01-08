package auth

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type UserService interface {
	SignUp(ctx context.Context, u model.User) error 
	SignIn(ctx context.Context, u model.User) (string, error)
	//ParseToken()
}