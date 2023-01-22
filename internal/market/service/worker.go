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

func fanOut(inputCh chan model.Accrual, stop <-chan bool, limitWorkers int) []chan model.Accrual {
	chs := make([]chan model.Accrual, 0, limitWorkers)
	for i := 0; i < limitWorkers; i++ {
		ch := make(chan model.Accrual)
		chs = append(chs, ch)
	}

	go func() {
		select {
		case <-stop:
			for i, ch := range chs {
				close(ch)
				log.Printf("закрыт in канал для воркера %d", i)
			}
			log.Print("каналы закрыты")
			return
		default:
			//log.Print("fanOut default")
			for i := 0; ; i++ {
				//log.Printf("Канал %d", i)
				if i == len(chs) {
					i = 0
				}

				order, ok := <-inputCh // если закрыт канал
				//log.Printf("чтение с общего канала %+v", order)
				if !ok {
					log.Print("горутина уснула")
					time.Sleep(30 * time.Second) //TODO как усыпить горутину?
					log.Print("горутина проснулась")
				}

				ch := chs[i]
				//log.Print("записываем и блокируемся")
				ch <- order
				//log.Printf("в канал #%d записали %v", i, order)
			}
		}
	}()

	return chs
}

func fanIn(inputChs ...chan model.Accrual) chan model.Accrual {
	outCh := make(chan model.Accrual)

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan model.Accrual) {
				defer wg.Done()
				for item := range inputCh { //закрыть канал для выхода из цикла
					outCh <- item
				}
			}(inputCh)
		}

		wg.Wait()
		close(outCh) //закрыть общий out - убрать
	}()

	return outCh
}

type AccrualResp struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func newWorker(in, out chan model.Accrual, i int, client *resty.Client) {
	go func() {
		for a := range in { //close chan
			log.Printf("worker #%d read accrual: %+v", i, a)

			var isProcessed bool
			for !isProcessed {
				URL := fmt.Sprintf("/api/orders/%s", a.Number)
				fmt.Println(URL)
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

				time.Sleep(3 * time.Second)
				continue

			}

			log.Printf("worker #%d stop working with accrual: %s", i, a.Number)

		}
	}()
}

func pollOrders(ctx context.Context, db market.AccrualRepository, out chan model.Accrual, pollInt time.Duration) error {
	t := time.NewTicker(pollInt)
	for {
		select {
		case <-ctx.Done():
			log.Printf("pollOrders is aborted by ctx")

			close(out)
			return nil

		case <-t.C:
			accruals, err := db.TakeOrdersForProcess(ctx, limitPollOrders)
			if err != nil {
                                log.Printf("")
				return err
			}

                        if len(accruals) == 0 {
                                log.Print("service doesn't have new orders for processing. Waiting...")
                                time.Sleep(5 * time.Second)
                                continue
                        }

			log.Printf("worker poll orders for processing: %v", accruals)

			err = db.ChangeStatusOrdersForProcess(ctx, accruals...)
			if err != nil {
				return err
			}

			for _, a := range accruals {
				out <- a
			}
		}
	}
}
