package main

import (
	"net/http"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/delivery/rest"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/delivery/rest/middleware"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/repository/localstorage"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/service"
	"github.com/labstack/echo/v4"
)

func main() {
	authStorage := localstorage.New()
	authService, err := service.New(authStorage)
	if err != nil {
		return
	}
	
	authHandler := rest.New(authService)

	e := echo.New()

	//routing
	e.GET("/", hello)

	e.POST("/api/user/register", authHandler.RegisterUser(authHandler.LoginUser))
	e.POST("/api/user/login", authHandler.LoginUser)

	g := e.Group("/api/user")
	g.Use(middleware.Auth)

	g.POST("/orders", postOrder)

	g.GET("/orders", getOrder)
	g.GET("/balance", getBalance)

	g.POST("/balance/withdraw", withdrawPoint)

	g.GET("/balance/withdrawals", getWithdrawal)

	e.Logger.Fatal(e.Start(":8080"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func postOrder(c echo.Context) error {
	return c.String(http.StatusOK, "postOrder")
}

func getOrder(c echo.Context) error {
	return c.String(http.StatusOK, "getOrder")
}

func getBalance(c echo.Context) error {
	return c.String(http.StatusOK, "getBalance")
}

func withdrawPoint(c echo.Context) error {
	return c.String(http.StatusOK, "withdrawPoint")
}

func getWithdrawal(c echo.Context) error {
	return c.String(http.StatusOK, "getWithdrawal")
}

/*
Накопительная система лояльности «Гофермарт» должна предоставлять следующие HTTP-хендлеры:

- `POST /api/user/register` — регистрация пользователя;
- `POST /api/user/login` — аутентификация пользователя;
- `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта;
- `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
- `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя;
- `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
- `GET /api/user/balance/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем.
*/
