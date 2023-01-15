package dto

import (
	"errors"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/jackc/pgtype"
)

type WithdrawnDTO struct {
	OrderID     string           `db:"order_id"`
	Sum         float64          `db:"sum"`
	ProcessedAt pgtype.Timestamp `db:"processed_at,omitempty"`
}

func WithdrawnToDTO(w model.Withdrawn) WithdrawnDTO {
	return WithdrawnDTO{
		OrderID: w.OrderID,
		Sum:     w.Sum,
	}
}

func WithdrawnsToModel(withdrawnsDTO ...WithdrawnDTO) ([]model.Withdrawn, error) {
	var withdrawns []model.Withdrawn

	for _, w := range withdrawnsDTO {
		val, err := w.ProcessedAt.Value()
		if err != nil {
			return nil, err
		}

		timeStamp, ok := val.(time.Time)
		if !ok {
			return nil, errors.New("error type assertion timestamp")
		}

		withdrawn := model.Withdrawn{
			OrderID:     w.OrderID,
			Sum:         w.Sum,
			ProcessedAt: timeStamp.Format(time.RFC3339),
		}

		withdrawns = append(withdrawns, withdrawn)

	}

	return withdrawns, nil
}
