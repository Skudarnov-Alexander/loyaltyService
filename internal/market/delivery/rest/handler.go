package rest

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/delivery/rest/dto"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service market.MarketService
}

func New(s market.MarketService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) PostOrder(c echo.Context) error {
	userID := c.Get("uuid") //TODO асерт типа
	/*
		if !ok {
			err := errors.New("uuid value is not string")
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	*/

	if userID == "" {
		err := errors.New("uuid value in context is empty")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("read body error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	orderID := string(data)

	log.Printf("orderID from body: %s", orderID)
	ctx := c.Request().Context()

	if err := h.service.SaveOrder(ctx, userID.(string), orderID); err != nil {
		log.Printf("service error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetOrders(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	ctx := c.Request().Context()

	userID := c.Get("uuid")

	if userID == "" {
		err := errors.New("uuid value in context is empty")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	orders, err := h.service.FetchOrders(ctx, userID.(string))
	if err != nil {
		err := fmt.Errorf("service error %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, orders)

}

func (h *Handler) GetBalance(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	ctx := c.Request().Context()

	userID := c.Get("uuid")

	if userID == "" {
		err := errors.New("uuid value in context is empty")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	balance, err := h.service.FetchBalance(ctx, userID.(string))
	if err != nil {
		err := fmt.Errorf("service error %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, balance)

}

func (h *Handler) GetWithdrawals(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	ctx := c.Request().Context()

	userID := c.Get("uuid")

	if userID == "" {
		err := errors.New("uuid value in context is empty")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	withdrawns, err := h.service.FetchWithdrawals(ctx, userID.(string))
	if err != nil {
		err := fmt.Errorf("service error %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

        withdrawnsDTO := dto.WithdrawnToDTO(withdrawns...)

	return c.JSON(http.StatusOK, withdrawnsDTO)

}

func (h *Handler) PostWithdrawal(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	ctx := c.Request().Context()

	userID := c.Get("uuid")

	if userID == "" {
		err := errors.New("uuid value in context is empty")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var wDTO dto.Withdrawn
	if err := c.Bind(&wDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

        w := dto.WithdrawnToModel(wDTO)

	b, err := h.service.MakeWithdrawal(ctx, userID.(string), w)
	if err != nil {
		log.Printf("service error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bDTO := dto.BalanceToDTO(b)

	return c.JSON(http.StatusOK, bDTO)

}
