package rest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/delivery/rest/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/pkg/luhn"

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

type Response struct {
        Message string `json:"message"`
}

func NewResponse(msg string) Response {
        return Response{
        	Message: msg,
        }
}

func (h *Handler) PostOrder(c echo.Context) error {
	userID, ok := c.Get("uuid").(string) //TODO асерт типа
        
        log.Printf("UUID from handler PostOrder: %s", userID)
	
	if !ok {
		err := errors.New("uuid value is not string")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	

	if userID == "" {
                log.Print("uuid value in context is empty")
		err := errors.New("uuid value in context is empty")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) //TODO error/log handler
	}

        ct := c.Request().Header.Get("Content-Type")
        if !strings.Contains(ct, "text/plain") {
                err := errors.New("header Content-Type is not text/plain")
                return echo.NewHTTPError(http.StatusBadRequest, err.Error())
        }

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("read body error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

        
        // TODO как првоерить пустое body?
        

	orderID := string(data)
        number, err := strconv.Atoi(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

        l := len(orderID)
	n := luhn.Checksum(number, l)
        fmt.Println(n)
	if n != 0 {
		//log.Printf("orderID is incorrect. Add %d to last num", div)
                return echo.NewHTTPError(http.StatusUnprocessableEntity, market.ErrFormatOrderID)
	}

	ctx := c.Request().Context()

        ok, err = h.service.CheckOrder(ctx, userID, orderID)
        if err != nil {
                log.Printf("service error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }

        if ok {
                return c.JSON(http.StatusOK, NewResponse("order is loaded yet"))
        }

	if err := h.service.SaveOrder(ctx, userID, orderID); err != nil {
                if errors.Is (err, market.ErrOrderIsExist) {
                        log.Printf("%v", err)
                        return echo.NewHTTPError(http.StatusConflict, err.Error())   
                }
		log.Printf("service error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

        
	return c.JSON(http.StatusAccepted, NewResponse("new order is loaded"))
}

/*
Возможные коды ответа:

- `200` — номер заказа уже был загружен этим пользователем;     +
- `202` — новый номер заказа принят в обработку;                +
- `400` — неверный формат запроса;                              +
- `401` — пользователь не аутентифицирован;                     +
- `409` — номер заказа уже был загружен другим пользователем;   +
- `422` — неверный формат номера заказа;                        +
- `500` — внутренняя ошибка сервера.                            +
*/

func (h *Handler) GetOrders(c echo.Context) error {
	userID := c.Get("uuid").(string) // TODO context interface my own
        log.Printf("UUID from handler GetOrders: %s", userID)

	if userID == "" {
		err := errors.New("uuid value in context is empty")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
        
        ctx := context.WithValue(c.Request().Context(), "uuid", c.Get("uuid"))  
        
	orders, err := h.service.FetchOrders(ctx, userID)
	if err != nil {
		err := fmt.Errorf("service error %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

        ordersDTO, err := dto.OrderToDTO(orders...)
        if err != nil {
                err := fmt.Errorf("parse time error %s", err.Error())
                return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }

        c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, ordersDTO)

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
