package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)




func VerifyJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		t, ok := c.Request().Header["Token"]

		fmt.Println(t)
		fmt.Println(ok)

		if !ok {
			return c.String(http.StatusBadRequest, "token header is empty")
		}

		token, err := jwt.Parse(t[0], func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				c.String(http.StatusUnauthorized, "Unauthorized")
				
			}



			return "", nil
		})

		if err != nil {
			c.String(http.StatusInternalServerError, "JWT parsing error")
		}

		if !token.Valid {
			c.String(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
		
		
	}
}