package dto

import "github.com/Skudarnov-Alexander/loyaltyService/internal/model"

type Withdrawn struct {
	OrderID     string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

func WithdrawnToDTO(withdrawns ...model.Withdrawn) []Withdrawn {
	var withdrawnsDTO []Withdrawn

	for _, w := range withdrawns {
		withdrawnDTO := Withdrawn{
			OrderID:     w.OrderID,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt,
		}

		withdrawnsDTO = append(withdrawnsDTO, withdrawnDTO)
	}

	return withdrawnsDTO

}

func WithdrawnToModel(w Withdrawn) model.Withdrawn {
	return model.Withdrawn{
		OrderID: w.OrderID,
		Sum:     w.Sum,
	}
}
