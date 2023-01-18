package service

import (
	"fmt"
	"log"
	"time"

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
		log.Print("fanOut go start")
		select {
		case <-stop:
			for i, ch := range chs {
                                close(ch)
                                log.Printf("закрыт in канал для воркера %d", i)
                        }
			log.Print("каналы закрыты")
			return
		default:
			log.Print("fanOut default")
			for i := 0; ; i++ {
				log.Printf("Канал %d", i)
				if i == len(chs) {
					i = 0
				}
	
				order, ok := <-inputCh // если закрыт канал
				log.Printf("чтение с общего канала %+v", order)
				if !ok {
					log.Print("горутина уснула")
					time.Sleep(30 * time.Second) //TODO как усыпить горутину?
					log.Print("горутина проснулась")
				}
	
				ch := chs[i]
				log.Print("записываем и блокируемся")
				ch <- order
				log.Printf("в канал #%d записали %v", i, order)
			}
		} 
    }()

    return chs
} 



func newWorker(in, out chan model.Accrual, i int, client *resty.Client) {
        go func() {
                for a := range in { //close chan
                        log.Printf("worker #%d read accrual: %+v", i, a)
                        var done bool
                        for !done {
                                log.Print("START LOOP")
                                URL := fmt.Sprintf("/api/orders/%s", a.Number)
                                resp, err := client.R().
                                SetResult(model.Accrual{}).
                                Get(URL)

                                if err != nil {
                                        log.Print(err)
                                        return
                                }

                                accrual := resp.Result().(*model.Accrual)
                                
                                log.Printf("Получен статус заказа по HTTP: %+v", accrual)
                                if accrual.Status == "REGISTERED" || accrual.Status == "PROCESSING" {
                                        time.Sleep(time.Minute)
                                        continue
                                }

                                out <- *accrual
                                done = true
                        } 
                        log.Printf("worker #%d stop working with accrual: %s", i, a.Number)
                }
        
            }()
}

