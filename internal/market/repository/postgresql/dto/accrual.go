package dto

import (
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/google/uuid"
)

type Accrual struct {
	Number  string    `db:"order_number"`
	Status  string    `db:"status"`
	Accrual float64   `db:"accrual"`
	UserID  uuid.UUID `db:"fk_user_id,omitempty"`
}

func AccrualToModel(a Accrual) model.Accrual {
	return model.Accrual{
		Number:  a.Number,
		Status:  a.Status,
		Accrual: a.Accrual,
		UserID:  a.UserID.String(),
	}
}
