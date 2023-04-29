package rest

import (
	"net/http"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/accrual/model"
	"github.com/labstack/echo/v4"
)

type accrualService interface {
	FetchAccrualStatus(order string) model.Accrual
	RegisterOrder(order model.Order) error
}

type Handler struct {
	service accrualService
}

func New(s accrualService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetAccrualStatus(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	number := c.Param("number")

	a := h.service.FetchAccrualStatus(number)
	return c.JSON(http.StatusOK, a)

}

func (h *Handler) PostNewOrder(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	var order model.Order
	if err := c.Bind(&order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.service.RegisterOrder(order); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
