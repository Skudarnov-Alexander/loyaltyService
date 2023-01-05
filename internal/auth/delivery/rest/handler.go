package rest

import (
	"fmt"
	"net/http"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service auth.UserService
}

func New(s auth.UserService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) RegisterUser(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	s := fmt.Sprintf("%+v", user)
	
	return c.String(http.StatusOK, s)
}

func (h *Handler) LoginUser(c echo.Context) error {
	return c.String(http.StatusOK, "loginUser")
}