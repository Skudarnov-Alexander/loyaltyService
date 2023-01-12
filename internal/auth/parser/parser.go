package parser

import (
	"errors"
	"fmt"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/service"
	"github.com/golang-jwt/jwt"
)

func ParseToken(accessToken string, singingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &service.Claims{}, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}

		return singingKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*service.Claims)
	if !ok {
		return "", errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt.Time.Unix() < time.Now().Unix() {
		return "", errors.New("token expired")
	}

	return claims.UserID, nil

}