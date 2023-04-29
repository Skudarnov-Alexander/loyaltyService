package dto

import (
	"errors"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Withdrawn struct {
	ID          int64            `db:"id,omitempty"`
	OrderNumber string           `db:"order_number"`
	Sum         float64          `db:"sum"`
	ProcessedAt pgtype.Timestamp `db:"processed_at,omitempty"`
	UserID      uuid.UUID        `db:"fk_user_id,omitempty"`
}

/*
func WithdrawnToDTO(w model.Withdrawn) WithdrawnDTO {
	return WithdrawnDTO{
		OrderID: w.OrderID,
		Sum:     w.Sum,
	}
}
*/

func (w *Withdrawn) ToModel() (*model.Withdrawn, error) {

	val, err := w.ProcessedAt.Value()
	if err != nil {
		return nil, err
	}

	timeStamp, ok := val.(time.Time)
	if !ok {
		return nil, errors.New("error type assertion timestamp")
	}

	return &model.Withdrawn{
		OrderID:     w.OrderNumber,
		Sum:         w.Sum,
		ProcessedAt: timeStamp.Format(time.RFC3339),
	}, nil

}
