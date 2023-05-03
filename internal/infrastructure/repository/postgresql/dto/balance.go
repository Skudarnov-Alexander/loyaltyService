package dto

import (
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/google/uuid"
)

type Balance struct {
	ID        int64     `db:"id,omitempty"`
	Current   float64   `db:"current_balance"`
	Withdrawn float64   `db:"withdrawn"`
	UserID    uuid.UUID `db:"fk_user_id,omitempty"`
}

func (b *Balance) ToModel() *model.Balance {
	return &model.Balance{
		Current:   b.Current,
		Withdrawn: b.Withdrawn,
	}
}
