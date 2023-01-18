package service

import (
	"log"
	"math/rand"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/accrual/model"
)

type accrualRepository interface {
        StoreNewOrder(order model.Order) error
}

type Service struct {
	db accrualRepository
}

func New(db accrualRepository) *Service{
	return &Service{
		db: db,
	}
}

func (s Service) FetchAccrualStatus(order string) model.Accrual {
	src := rand.NewSource(time.Now().UnixNano())
        r := rand.New(src)
        n := r.Intn(100)

        var status string
        var accrual float64

        switch {
        case n < 10:
                status = "REGISTERED"
        case n < 15:
                status = "INVALID"
        case n < 50:
                status = "PROCESSING"
        case n < 100:
                status = "PROCESSED"
                accrual = float64(r.Intn(5000))
        }

        a := model.Accrual{
        	Order:   order,
        	Status:  status,
        	Accrual: accrual,
        }

        log.Print(a)

        return a
}

func (s Service) RegisterOrder(order model.Order) error {
        if err := s.db.StoreNewOrder(order); err != nil {
                return err
        }

        return nil
}

/*
func requestToStatusOrder(o model.Accrual, st AccrualStorage) (AccrualStatus, error) {
        aStatus, ok := st.Storage[o.Number]
        if !ok {
                s := rand.NewSource(time.Now().UnixNano())
                r := rand.New(s)
                n := r.Intn(100)

                var status string

                status = "REGISTERED"
                if n < 20 {
                        status = "INVALID"
                }

                as := AccrualStatus{
                	Order:   o.Number,
                	Status:  status,
                }

                st.Storage[o.Number] = as

                log.Printf("заказ #%s записан в обработку со статусом %s", o.Number, status)

                go func(st AccrualStorage, orderNum string, n int) {
                        time.Sleep(time.Duration(n) * time.Second)

                        var status string
                        status = "PROCESSING"
                        st.Storage[orderNum]= AccrualStatus{
                                Order:   orderNum,
                                Status:  status,
                        }

                        log.Printf("заказ #%s записан в обработку со статусом %s", orderNum, status)

                        time.Sleep(time.Duration(n) * time.Second)
                        st.Storage[orderNum] = AccrualStatus{
                                Order:   orderNum,
                                Status:  status,
                        }


                        time.Sleep(time.Duration(n) * time.Second)
                        s := rand.NewSource(time.Now().UnixNano())
                        r := rand.New(s)

                        accrual := r.Intn(2000)
                        status = "PROCESSED"

                        st.Storage[orderNum] = AccrualStatus{
                                Order:   orderNum,
                                Status:  status,
                                Accrual: float64(accrual),
                        }

                }(st, o.Number, n)

                return as, nil
        }
        return aStatus, nil

}
*/