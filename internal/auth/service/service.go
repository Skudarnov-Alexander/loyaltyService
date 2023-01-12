package service

import (
	"context"
	"errors"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo auth.AuthRepository
	hashSalt []byte
}

func New(userRepo auth.AuthRepository) (*AuthService, error) {
	hashSalt, err := generateRandomSalt(saltSize)
	if err != nil {
		return nil, err
	}
	
	return &AuthService{
		userRepo: userRepo,
		hashSalt: hashSalt,
	}, nil
}

func (s *AuthService) SignUp(ctx context.Context, u model.User) error {
	u.Password = hashPassword(u.Password, s.hashSalt)
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	
	u.ID = uuid.String()

	return s.userRepo.CreateUser(ctx, u)
}

func (s *AuthService) SignIn(ctx context.Context, u model.User) (string, error) {
	user, err := s.userRepo.GetUser(ctx, u.Username)
	if err != nil {
		return "", err
	}
	
	if ok := doPasswordsMatch(user.Password, u.Password, s.hashSalt); !ok {
		return "", errors.New("pass is not match")
	}

	return generateJWT(user)

}



