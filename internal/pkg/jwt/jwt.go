package jwt

import (
	"fmt"
	"time"

	error2 "github.com/Skudarnov-Alexander/loyaltyService/internal/error"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Username string
	UserID   string
}

func GenerateJWT(key []byte, u model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(10 * time.Minute),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		Username: u.Username,
		UserID:   u.ID,
	})

	return token.SignedString(key)

}

func ParseToken(accessToken string, singingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}

		return singingKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", error2.ErrInvalidAccessToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", error2.ErrClaimsParsing
	}

	if claims.ExpiresAt.Time.Unix() < time.Now().Unix() {
		return "", error2.ErrExpiredToken
	}

	return claims.UserID, nil

}
