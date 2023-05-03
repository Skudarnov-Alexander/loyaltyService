package main

import (
	"github.com/Skudarnov-Alexander/loyaltyService/internal/accrual/delivery/rest"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/accrual/service"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/accrual/repository/storage"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	storage := storage.New()
	service := service.New(storage)
	handler := rest.New(service)

	e.GET("/api/orders/:number", handler.GetAccrualStatus)
	e.POST("/api/orders", handler.PostNewOrder)

	e.Logger.Fatal(e.Start(":8082"))
}

