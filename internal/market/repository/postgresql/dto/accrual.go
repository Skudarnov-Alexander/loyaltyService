package dto

import "github.com/Skudarnov-Alexander/loyaltyService/internal/model"

type Accrual struct {
	Number  string   `db:"order_number"`
    Status  string   `db:"status"`
    Accrual float64  `db:"accrual"`
}

func AccrualToModel(a Accrual) model.Accrual {
	return model.Accrual{
		Number:  a.Number,
		Status:  a.Status,
		Accrual: a.Accrual,
	}
}
