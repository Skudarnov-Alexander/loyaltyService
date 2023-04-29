package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

/*
type response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg,omitempty"`
	Token  string `json:"token,omitempty"`
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

*/

type authInteractor interface {
	SignUp(ctx context.Context, username, pwd string) error
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

/*
type User struct {
	ID       string
	Username string `json:"login" validate:"required,ascii"`
	Password string `json:"password" validate:"required,ascii"`
    }
*/
func (ac *AuthHTTPController) SignUp(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	reqData := new(SignUpRequest)
	if err := c.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ac.validator.Struct(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	if err := ac.authInteractor.SignUp(ctx, reqData.Username, reqData.Password); err != nil {
		if errors.Is(err, auth.ErrUserIsExist) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "new user is created")
	/*
		data, err := json.Marshal(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		b := bytes.NewBuffer(data)

		body := io.NopCloser(b)

		c.Request().Body = body

		return next(c)
	*/
}

/*
Возможные коды ответа:

- `200` — пользователь успешно зарегистрирован и аутентифицирован;
- `400` — неверный формат запроса;
- `409` — логин уже занят;
- `500` — внутренняя ошибка сервера.
*/
/*
func (ac *AuthHTTPController) LogIn(c echo.Context) error {
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

	c.Response().Header().Add("Authorization", token)

	return c.JSON(http.StatusOK, newResponse(http.StatusOK, "auth is successfull", token))
}

/*
func (h *Handler) RegisterUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		var user model.User

		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
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


Возможные коды ответа:

- `200` — пользователь успешно зарегистрирован и аутентифицирован;
- `400` — неверный формат запроса;
- `409` — логин уже занят;
- `500` — внутренняя ошибка сервера.


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

	c.Response().Header().Add("Authorization", token)

	return c.JSON(http.StatusOK, newResponse(http.StatusOK, "auth is successfull", token))
}

*/
