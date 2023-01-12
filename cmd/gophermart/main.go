package main

import (
	"log"

	authr "github.com/Skudarnov-Alexander/loyaltyService/internal/auth/delivery/rest"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/delivery/rest/middleware"
	authdb "github.com/Skudarnov-Alexander/loyaltyService/internal/auth/repository/postgresql"
	auths "github.com/Skudarnov-Alexander/loyaltyService/internal/auth/service"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/database"
	marketr "github.com/Skudarnov-Alexander/loyaltyService/internal/market/delivery/rest"
	marketdb "github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql"
	markets "github.com/Skudarnov-Alexander/loyaltyService/internal/market/service"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.InitDB(db); err != nil {
		log.Fatal(err)
	}

	authStorage, err := authdb.New(db)
	if err != nil {
		log.Fatal()
	}

	authService, err := auths.New(authStorage)
	if err != nil {
		return
	}

	authHandler := authr.New(authService)

	marketStorage, err := marketdb.New(db)
	if err != nil {
		log.Fatal()
	}

	marketService := markets.New(marketStorage)
	if err != nil {
		return
	}

	marketHandler := marketr.New(marketService)

	e := echo.New()

	e.POST("/api/user/register", authHandler.RegisterUser(authHandler.LoginUser))
	e.POST("/api/user/login", authHandler.LoginUser)

	g := e.Group("/api/user")
	g.Use(middleware.Auth)

	g.POST("/orders", marketHandler.PostOrder)
	g.GET("/orders", marketHandler.GetOrders)
	g.GET("/balance", marketHandler.GetBalance)
	g.POST("/balance/withdraw", marketHandler.PostWithdrawal)
	g.GET("/balance/withdrawals", marketHandler.GetWithdrawals)

	e.Logger.Fatal(e.Start(":8080"))
}
