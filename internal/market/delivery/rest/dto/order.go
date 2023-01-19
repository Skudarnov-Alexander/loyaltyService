package dto

import "github.com/Skudarnov-Alexander/loyaltyService/internal/model"

type Order struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

func OrderToDTO(orders ...model.Order) []Order {
	var ordersDTO []Order

	for _, o := range orders {
		orderDTO := Order{
			Number:     o.Number,
			Status:     o.Status,
			UploadedAt: o.UploadedAt,
		}

		if o.Status == "PROCESSED" {
			orderDTO.Accrual = &o.Accrual
		}

		ordersDTO = append(ordersDTO, orderDTO)

	}

	return ordersDTO
}
