package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/config"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/database"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/delivery/rest"
	http2 "github.com/Skudarnov-Alexander/loyaltyService/internal/infrastructure/http"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/infrastructure/repository/postgresql"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/logger"
	marketr "github.com/Skudarnov-Alexander/loyaltyService/internal/market/delivery/rest"
	marketdb "github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql"
	markets "github.com/Skudarnov-Alexander/loyaltyService/internal/market/service"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/pkg/hash"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/server"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/usecase/interactor"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), 
                syscall.SIGINT, 
                syscall.SIGTERM, 
                syscall.SIGQUIT)

	defer cancel()

        // init logger
        logger := logger.New()

        // create error channel
	errChan := make(chan error)
	go func() {
		for err := range errChan {
			logger.L.Err(err).Msg("ErrChan")
		}
	}()

        // parse config
	cfg, err := config.New()
	if err != nil {
		logger.L.Fatal().Err(err).Msg("config parsing error")
	}
        logger.L.Info().Msgf("Config: %+v\n", cfg)
       
        // init DB connection and create tables/test data
	db, err := database.New(cfg.DBAddr)
	if err != nil {
		logger.L.Fatal().Err(err).Msg("DB init error")
	}
        logger.L.Info().Msgf("DB is connected: %s\n", cfg.DBAddr)

	if err := database.CreateTables(db); err != nil {
		logger.L.Fatal().Err(err).Msg("DB tables creation error")
	}

        // init auth service
	/*
	aStorage := authdb.New(db)
	aService, err := auths.New(aStorage) //TODO убрать соль
	if err != nil {
		log.Fatal(err)
	}

	aHandler := authr.New(aService)
	*/

        // init gophermarket service and accrual worker
	mStorage := marketdb.New(db)
	mService := markets.New(mStorage)
	mHandler := marketr.New(mService)

        accrualService := markets.NewAccrualService(mStorage, cfg.PollInt, errChan)

        // init HTTP server
	server := server.New(nil, mHandler, cfg.Addr)

	
        // start App
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.Run()
	})

	g.Go(func() error {
		return accrualService.Run(ctx, cfg.AccrualAddr)
	})

	salt, err := hash.GenerateRandomSalt()
	if err != nil {
		logger.L.Fatal().Err(err).Msg("salt generating error")
	}
	hasher := hash.New(salt)

	userRepository := postgresql.NewUserRepository(db)
	balanceRepository := postgresql.NewBalanceRepository(db)

	authInteractor := interactor.NewAuthInteractor(userRepository, balanceRepository, hasher)
	authHTTPConroller := rest.NewAuthHTTPController(authInteractor)

	echoServer := http2.NewEchoHTTPServer(authHTTPConroller, nil)
	echoServer.Run(8086)

        // gracefull server shutdown
        go func() {
                <-ctx.Done()
                server.Stop(ctx)
        }()
        
	if err = g.Wait(); err != nil {
		log.Print(err)
	}

        // App gracefull shutdown
        log.Print("App is shutting down...")
	time.Sleep(10 * time.Second) //Q какие ресурсы надо закрыть?
        defer db.Close()
	log.Print("Agent is shutted down")
}
