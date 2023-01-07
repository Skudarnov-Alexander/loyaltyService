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
	c.Response().Header().Set("Content-Type", "application/json")
	user := &model.User{}
	
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, `{"error":"bind error"}`)
	}

	err := h.service.SignUp(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, `{"error":"service error"}`)
	}

	return c.JSON(http.StatusOK, user)
}

type response struct {
	Status int	`json:"status"`
	Msg	string	`json:"msg,omitempty"`
	Token string `json:"token,omitempty"`
}

func newResponse(status int, msg, token string) *response {
	return &response{
		Status: status,
		Msg:    msg,
		Token:  token,
	}
}
func (h *Handler) LoginUser(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	user := &model.User{}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, `{"error":"bind error"}`)
	}
	token, err := h.service.SignIn(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, `{"error":"service error"}`)
	}
	fmt.Printf("token: %s", token)
	return c.JSON(http.StatusOK, newResponse(http.StatusOK, "auth is successfull", token))
}



