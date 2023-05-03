package http

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthHTTPController interface {
	HandleSignUp(c echo.Context) error
	HandleLogIn(c echo.Context) error
}

type MarketHTTPController interface {
}

type EchoHTTPServer struct {
	e          *echo.Echo
	authCtrl   AuthHTTPController
	marketCtrl MarketHTTPController
}

func NewEchoHTTPServer(
	authCtrl AuthHTTPController,
	marketCtrl MarketHTTPController,
) *EchoHTTPServer {
	server := &EchoHTTPServer{
		e:          echo.New(),
		authCtrl:   authCtrl,
		marketCtrl: marketCtrl,
	}

	server.registerAuthGroup()
	server.registerMarketGroup()

	return server
}

func (s *EchoHTTPServer) registerAuthGroup() {
	authGroup := s.e.Group("/api/user")

	authGroup.POST("/register", s.handleRegisterUser)
	authGroup.POST("/login", s.handleLogInUser)
}

func (s *EchoHTTPServer) registerMarketGroup() {
	marketGroup := s.e.Group("/api/user")

	_ = marketGroup

}

func (s *EchoHTTPServer) handleRegisterUser(c echo.Context) error {
	return s.authCtrl.HandleSignUp(c)
}

func (s *EchoHTTPServer) handleLogInUser(c echo.Context) error {
	return s.authCtrl.HandleLogIn(c)
}

func (s *EchoHTTPServer) Run(port int) {
	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Printf("HTTP Server is starting: %s", addr)

		if err := s.e.Start(addr); err != nil {
			log.Error(err)
		}
	}()
}

func (s *EchoHTTPServer) Shutdown(ctx context.Context) error {
	log.Printf("HTTP Server is stopping: %s", s.e.Server.Addr)

	return s.e.Shutdown(ctx)
}
