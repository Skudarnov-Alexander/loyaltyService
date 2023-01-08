package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/parser"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/service"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			err := errors.New("header Authorization is empty")
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			err := errors.New("header Authorization is incorrect")
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		if headerParts[0] != "Bearer" {
			err := errors.New("header Authorization is not Bearer")
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}


		err := parser.ParseToken(headerParts[1], service.SampleSecretKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return next(c)
		
	}
}

