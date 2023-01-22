package main

import (
	"context"
	"log"

	authr "github.com/Skudarnov-Alexander/loyaltyService/internal/auth/delivery/rest"
	authdb "github.com/Skudarnov-Alexander/loyaltyService/internal/auth/repository/postgresql"
	auths "github.com/Skudarnov-Alexander/loyaltyService/internal/auth/service"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/config"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/database"
	marketr "github.com/Skudarnov-Alexander/loyaltyService/internal/market/delivery/rest"
	marketdb "github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql"
	markets "github.com/Skudarnov-Alexander/loyaltyService/internal/market/service"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/server"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	go func() {
		for err := range errChan {
			log.Print(err)
		}
	}()

	cfg, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	db, err := database.New(cfg.DBAddr)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	aStorage, err := authdb.New(db)
	if err != nil {
		log.Fatal(err)
	}

	aService, err := auths.New(aStorage)
	if err != nil {
		log.Fatal(err)
	}

	aHandler := authr.New(aService)

	mStorage, err := marketdb.New(db)
	if err != nil {
		log.Fatal()
	}

	mService := markets.New(mStorage)
	if err != nil {
		log.Fatal(err)
	}

	mHandler := marketr.New(mService)

	server := server.New(aHandler, mHandler, cfg.Addr)

	accrualService := markets.NewAccrualService(mStorage, cfg.PollInt)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.Run()
	})

	g.Go(func() error {
		return accrualService.Run(ctx, cfg.AccrualAddr)
	})

	if err = g.Wait(); err != nil {
		cancel()
		log.Fatal(err)
	}

}
