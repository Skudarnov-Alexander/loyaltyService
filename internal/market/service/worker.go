package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/go-resty/resty/v2"
)

type AccrualResp struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func fanOut(inputCh chan model.Accrual, limitWorkers int) []chan model.Accrual {
	chs := make([]chan model.Accrual, 0, limitWorkers)
	for i := 0; i < limitWorkers; i++ {
		ch := make(chan model.Accrual)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan model.Accrual) {
			for i, ch := range chs {
				close(ch)
				log.Printf("закрыт in канал для воркера %d", i)
			}
			log.Print("Все каналы in закрыты")
		}(chs)

		for i := 0; ; i++ {

			order, ok := <-inputCh
			if !ok {
				log.Print("fanOut is stopped by сlosing a input chan")
				return
			}

			//log.Printf("Канал %d", i)
			if i == len(chs) {
				i = 0
			}

			ch := chs[i]
			//log.Print("записываем и блокируемся")
			ch <- order
			//log.Printf("в канал #%d записали %v", i, order)

		}

	}()

	return chs
}

func fanIn(inChs ...chan model.Accrual) chan model.Accrual {
	log.Print("fanIn is started")

	out := make(chan model.Accrual)

	go func() {
		wg := &sync.WaitGroup{}

		for _, inCh := range inChs {
			wg.Add(1)

			go func(inCh chan model.Accrual) {
				defer wg.Done()
				for accrual := range inCh { //закрыть канал out воркера для выхода из цикла со стороны воркера
					out <- accrual
				}
			}(inCh)
		}

		wg.Wait()
                close(out) //закрыть общий out - убрать
	}()

	
	log.Print("fanIn is stopped")

	return out
}

func newWorker(in <-chan model.Accrual, out chan<- model.Accrual, i int, client *resty.Client) {
	go func() {
		for a := range in { //close chan
			log.Printf("worker #%d read accrual: %+v", i, a)

			URL := fmt.Sprintf("/api/orders/%s", a.Number)
			fmt.Println(URL)

			var isProcessed bool
			for !isProcessed {
				resp, err := client.R().
					SetResult(AccrualResp{}).
					Get(URL)

				if err != nil {
					log.Print(err)
					return
				}

				accrualResp := resp.Result().(*AccrualResp)

				log.Printf("Получен статус заказа по HTTP: %+v", accrualResp)

				if accrualResp.Status == "PROCESSED" || accrualResp.Status == "INVALID" {
					accrual := model.Accrual{
						Number:  accrualResp.Number,
						Status:  accrualResp.Status,
						Accrual: accrualResp.Accrual,
						UserID:  a.UserID,
					}

					out <- accrual
					isProcessed = true

				}

				time.Sleep(3 * time.Second) //retry to request to accrual API
				continue

			}
		}
                log.Printf("worker #%d finished working with accrual", i)
		close(out)
	}()
}

func pollOrders(accruals ...model.Accrual) chan model.Accrual {
	out := make(chan model.Accrual)

	go func(out chan model.Accrual) {
		for _, a := range accruals {
			out <- a
		}

		close(out)
	}(out)

	return out
}

func writeResult(ctx context.Context, db market.AccrualRepository, in chan model.Accrual) {
	go func() {
		for accrual := range in {
			if err := db.UpdateStatusProcessedOrders(accrual); err != nil {
				log.Printf("readWorker error: %s", err)
				//return err
			}

			if accrual.Status == "PROCESSED" {
				if err := db.UpdateBalanceProcessedOrders(accrual); err != nil {
					log.Printf("readWorker error: %s", err)
					//return err
				}
			}
		}

                log.Print("writeResult finished")
	}()
}
