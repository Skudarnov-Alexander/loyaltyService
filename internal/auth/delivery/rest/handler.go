package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/labstack/echo/v4"
)

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

type Handler struct {
	service auth.UserService
}

func New(s auth.UserService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) RegisterUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	var user model.User
	
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.service.SignUp(c.Request().Context(), user)
	if errors.Is(err, auth.ErrUserIsExist) {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	} 
	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data, err := json.Marshal(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	b := bytes.NewBuffer(data)

	body := io.NopCloser(b)

	c.Request().Body = body

	return next(c)
	
	}
}

/*
Возможные коды ответа:

- `200` — пользователь успешно зарегистрирован и аутентифицирован;
- `400` — неверный формат запроса;
- `409` — логин уже занят;
- `500` — внутренняя ошибка сервера.
*/


func (h *Handler) LoginUser(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	var user model.User

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := h.service.SignIn(c.Request().Context(), user)
	if errors.Is(err, auth.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	} 
	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, newResponse(http.StatusOK, "auth is successfull", token))
}

/*
Возможные коды ответа:

- `200` — пользователь успешно аутентифицирован;
- `400` — неверный формат запроса;
- `401` — неверная пара логин/пароль;
- `500` — внутренняя ошибка сервера.
*/



