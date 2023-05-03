package rest

import (
	"context"
	"errors"
	"net/http"

	error2 "github.com/Skudarnov-Alexander/loyaltyService/internal/error"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type authInteractor interface {
	SignUp(ctx context.Context, username, pwd string) (string, error)
	LogIn(ctx context.Context, username, pwd string) (string, error)
}

type AuthHTTPController struct {
	authInteractor
	validator *validator.Validate
}

func NewAuthHTTPController(authInteractor authInteractor) *AuthHTTPController {
	return &AuthHTTPController{
		authInteractor: authInteractor,
		validator:      validator.New(),
	}
}

type SignUpRequest struct {
	Username string `json:"login" validate:"required,ascii"`
	Password string `json:"password" validate:"required,ascii"`
}

type LogInRequest struct {
	Username string `json:"login" validate:"required,ascii"`
	Password string `json:"password" validate:"required,ascii"`
}

type AuthMarket struct {
	JWT string `json:"jwt"`
}

// Возможные коды ответа:
//
// 200 — пользователь успешно зарегистрирован и аутентифицирован;
// 400 — неверный формат запроса;
// 409 — логин уже занят;
// 500 — внутренняя ошибка сервера.
func (ac *AuthHTTPController) HandleSignUp(c echo.Context) error {
	reqData := new(SignUpRequest)

	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ac.validator.Struct(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	token, err := ac.authInteractor.SignUp(ctx, reqData.Username, reqData.Password)
	if err != nil {
		if errors.Is(err, error2.ErrUserIsExist) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Add("Authorization", token)

	return c.JSON(http.StatusOK, AuthMarket{
		JWT: token,
	})
}

// Возможные коды ответа:
//
// 200 — пользователь успешно аутентифицирован;
// 400 — неверный формат запроса;
// 401 — неверная пара логин/пароль;
// 500 — внутренняя ошибка сервера.
func (ac *AuthHTTPController) HandleLogIn(c echo.Context) error {
	reqData := new(LogInRequest)

	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ac.validator.Struct(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	token, err := ac.authInteractor.LogIn(ctx, reqData.Username, reqData.Password)
	if err != nil {
		if errors.Is(err, error2.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}

	c.Response().Header().Add("Authorization", token)

	return c.String(http.StatusOK, token)
}
