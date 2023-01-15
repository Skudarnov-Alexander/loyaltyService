package dto

import "github.com/Skudarnov-Alexander/loyaltyService/internal/model"

type BalanceDTO struct {
	Current   float64 `db:"current_balance"`
	Withdrawn float64 `db:"withdrawn"`
}

func BalanceToModel(b BalanceDTO) model.Balance {
	return model.Balance{
		Current:   b.Current,
		Withdrawn: b.Withdrawn,
	}
}
