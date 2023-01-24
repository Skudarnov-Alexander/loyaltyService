package storage

import (
	"errors"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/accrual/model"
)

type Storage struct {
	st map[string]model.Order
} 

func New() Storage {
	return Storage{
		st: make(map[string]model.Order),
	}
}

func (st Storage) StoreNewOrder(order model.Order) error {
	_, ok := st.st[order.Order]
	if ok {
		return errors.New("order is loaded yet")
	}

	status := "REGISTERED"
	
	order.Status = status

	st.st[order.Order] = order

	log.Printf("заказ #%+v для обработки загружен", order)
	return nil
}

