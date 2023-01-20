package dto

import (
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

type Order struct {
	Number     string       `json:"number"`
	Status     string       `json:"status"`
	Accrual    float32      `json:"accrual,omitempty"`
	UploadedAt time.Time    `json:"uploaded_at"`
}
/*
type order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}
*/

func OrderToDTO(orders ...model.Order) ([]Order, error) {
	var ordersDTO []Order

	for _, o := range orders {
                timeStamp, err := time.Parse(time.RFC3339, o.UploadedAt)
                if err != nil {
                        return nil, err
                }
		orderDTO := Order{
			Number:     o.Number,
			Status:     o.Status,
			Accrual:    float32(o.Accrual),
			UploadedAt: timeStamp,
		}
                /*
		if o.Status == "PROCESSED" {
			orderDTO.Accrual = &o.Accrual
		}
                */

		ordersDTO = append(ordersDTO, orderDTO)

	}

	return ordersDTO, nil
}
