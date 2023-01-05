package service

import "github.com/Skudarnov-Alexander/loyaltyService/internal/auth"

type AuthService struct {
	userRepo  auth.UserRepository
}

func New(userRepo auth.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) SignUp() {}

func (s *AuthService) SignIn() {}