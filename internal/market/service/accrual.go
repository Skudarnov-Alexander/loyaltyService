package service

import (
	"context"
	"log"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/go-resty/resty/v2"
)

const (
	limitWorkers    int = 5
	limitPollOrders int = 10
)

type AccrualService struct {
	db      market.AccrualRepository
	pollInt time.Duration
        errChan chan<- error
}

func NewAccrualService(db market.AccrualRepository, pollInt time.Duration, errChan chan<- error) *AccrualService {
	return &AccrualService{
		db:      db,
		pollInt: pollInt,
	}
}

func (s AccrualService) Run(ctx context.Context, accrualAddr string) error {
	log.Print("Accrual service is started")
	client := resty.New().
		SetBaseURL(accrualAddr)

	t := time.NewTicker(s.pollInt)
        var count int

        defer log.Print("Accrual service is aborted by ctx") //TODO грейсфул
	for {
		select {
		case <-ctx.Done():
                        return nil

		case <-t.C:
			accruals, err := s.db.TakeOrdersForProcess(ctx, limitPollOrders)
			if err != nil {
				log.Print(err)
				return err
			}

                        log.Printf("worker poll orders for processing: %v", accruals)

			
			if len(accruals) == 0 {
				count++

				if count == 5 {
					log.Print("service doesn't have new orders for processing a lot of time. Waiting...")
                                        count = 0
					time.Sleep(15 * time.Second) //waiting new orders in DB
				}

				continue
			}

			if err := s.db.ChangeStatusOrdersForProcess(ctx, accruals...); err != nil {
				log.Print(err)
				return err
			}
                        
                        // TODO надо ли вейт групп на эти горутины все?
			inAccrualCh := pollOrders(accruals...)
			fanOutChs := fanOut(inAccrualCh, limitWorkers)

			workerOutChs := make([]chan model.Accrual, 0, limitWorkers)
			for idx, fanOutCh := range fanOutChs {
				workerOutCh := make(chan model.Accrual)
				newWorker(fanOutCh, workerOutCh, idx, client)
				workerOutChs = append(workerOutChs, workerOutCh)
			}

			outAccrualCh := fanIn(workerOutChs...)

			writeResult(ctx, s.db, outAccrualCh)

		}
	}
}
