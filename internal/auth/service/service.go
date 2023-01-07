package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/golang-jwt/jwt/v4"
)

var sampleSecretKey = []byte("SecretYouShouldHide")


type AuthService struct {
	userRepo auth.UserRepository
	hashSalt []byte
}

func New(userRepo auth.UserRepository) (*AuthService, error) {
	hashSalt, err := generateRandomSalt(saltSize)
	if err != nil {
		return nil, err
	}
	
	return &AuthService{
		userRepo: userRepo,
		hashSalt: hashSalt,
	}, nil
}

func (s *AuthService) SignUp(ctx context.Context, u *model.User) error {
	u.Password = hashPassword(u.Password, s.hashSalt)
	return s.userRepo.CreateUser(ctx, u)
}

func (s *AuthService) SignIn(ctx context.Context, u *model.User) (string, error) {
	user, err := s.userRepo.GetUser(ctx, u.Username)
	if err != nil {
		return "", err
	}
	
	if ok := doPasswordsMatch(user.Password, u.Password, s.hashSalt); !ok {
		fmt.Println("not OK")
		return "", errors.New("pass is not match")
	}

	fmt.Println("before JWT")
	return generateJWT(u.Username)

}

type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}



func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(2 * time.Minute),
			},
			IssuedAt:  &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		Username:         username,
	})
	
	str, err := token.SignedString(sampleSecretKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	
	return  str, nil
	
}
