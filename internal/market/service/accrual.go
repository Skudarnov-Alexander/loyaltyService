package service

import (
	"context"
	"log"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/go-resty/resty/v2"
)

type AccrualService struct {
	db market.AccrualRepository
	pollInt time.Duration
}

func NewAccrualService(db market.AccrualRepository, pollInt time.Duration) *AccrualService {
	return &AccrualService{
		db:      db,
		pollInt: pollInt,
	}
}

func (s AccrualService) Run(ctx context.Context, accrualAddr string) error {
	client := resty.New().
		SetBaseURL("http://127.0.0.1:8082")
	

	var count int64
	inAccrualCh := make(chan model.Accrual, 100)
	outAccrualCh:= make(chan model.Accrual, 100)
	stop := make(chan bool)

	fanOutChs := fanOut(inAccrualCh, stop, 3)

	for idx, fanOutCh := range fanOutChs {
		newWorker(fanOutCh, outAccrualCh, idx, client)
	}

       go func() {
              if err := readWorker(ctx, s.db, outAccrualCh); err != nil {
                     log.Print(err)
                     return
              }
       }() 
	
	t := time.NewTicker(s.pollInt)
	for {
		<- t.C
		count++

		accruals, err := s.db.TakeOrdersForProcess(ctx)
		if err != nil {
			return err
		}

		log.Printf("worker load orders: %v", accruals)

		err = s.db.ChangeStatusOrdersForProcess(ctx, accruals...)
		if err != nil {
			return err
		}

		for _, a := range accruals{
			inAccrualCh <- a
			log.Printf("Записали в обший канал %+v", a)
		}

		if count == 5 {
			break
		}
	}

	stop <- true

	time.Sleep(5 * time.Second)
	log.Print("AccrualService stopped")
	return nil
}

func readWorker(ctx context.Context, db market.AccrualRepository, in chan model.Accrual) error {
       for accrual := range in {
              if err := db.UpdateStatusProcessedOrders(ctx, accrual); err != nil {
                     log.Printf("readWorker error: %s", err)
                     return err
              }
       }

       return nil
}




