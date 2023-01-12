package service

import (
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

var SampleSecretKey = []byte("SecretYouShouldHide")

type Claims struct {
	jwt.RegisteredClaims
	Username string 
	UserID	 string
}

func generateJWT(u model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(2 * time.Minute),
			},
			IssuedAt:  &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		Username:         u.Username,
		UserID:           u.ID,
	})
	
	str, err := token.SignedString(SampleSecretKey)
	if err != nil {
		return "", err
	}
	
	return  str, nil
	
}

