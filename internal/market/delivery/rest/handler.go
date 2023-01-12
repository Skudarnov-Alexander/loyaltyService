package rest

import (
	"log"
	"net/http"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
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

type keyID string

func setKey (c echo.Context, key keyID) {
	c.Set("uuid", key)
}

func getKey (c echo.Context) (keyID, bool) {
	val, ok := c.Get("uuid").(keyID)
	return val, ok
	
}


func (h *Handler) PostOrder(c echo.Context) error {
		return c.String(http.StatusOK, "PostOrder")
}

func (h *Handler) GetOrders(c echo.Context) error {
	uuid, ok := getKey(c)
	if !ok {
		log.Fatal("keyID in context reading error")
	}
	return c.String(http.StatusOK, string(uuid))
	
}

func (h *Handler) GetBalance(c echo.Context) error {
	return c.String(http.StatusOK, "GetBalance")
	
}

func (h *Handler) GetWithdrawals(c echo.Context) error {
	return c.String(http.StatusOK, "GetWithdrawals")
	
}

func (h *Handler) PostWithdrawal(c echo.Context) error {
	return c.String(http.StatusOK, "PostWithdrawal")
	
}








