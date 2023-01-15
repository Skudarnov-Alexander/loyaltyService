package dto

import "github.com/Skudarnov-Alexander/loyaltyService/internal/model"

type Balance struct {
	CurrentUpdated   float64 `json:"current"`
	WithdrawnUpdated float64 `json:"withdrawn"`
}

func BalanceToDTO(b model.Balance) Balance {
	return Balance{
		CurrentUpdated:   b.Current,
		WithdrawnUpdated: b.Withdrawn,
	}
}
