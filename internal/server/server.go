package server

import (
	"log"

	arest "github.com/Skudarnov-Alexander/loyaltyService/internal/auth/delivery/rest"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/delivery/rest/middleware"
	mrest "github.com/Skudarnov-Alexander/loyaltyService/internal/market/delivery/rest"

	"github.com/labstack/echo/v4"
)



type Server struct {
	e *echo.Echo
}

func New(authHandler *arest.Handler, marketHandler *mrest.Handler, addr string) Server {
	e := echo.New()
        e.Server.Addr = addr

	e.POST("/api/user/register", authHandler.RegisterUser(authHandler.LoginUser))
	e.POST("/api/user/login", authHandler.LoginUser)

	g := e.Group("/api/user")
	g.Use(middleware.Auth)

	g.POST("/orders", marketHandler.PostOrder)
	g.GET("/orders", marketHandler.GetOrders)
	g.GET("/balance", marketHandler.GetBalance)
	g.POST("/balance/withdraw", marketHandler.PostWithdrawal)
	g.GET("/withdrawals", marketHandler.GetWithdrawals)

	return Server{
		e: e,
	}
}

func (s Server) Run() error {
        log.Printf("HTTP Server is starting: %s", s.e.Server.Addr)
	if err := s.e.Server.ListenAndServe(); err != nil {
                log.Printf("HTTP Server is stopping: %s", s.e.Server.Addr)   
                return err
        }

        return nil

}
