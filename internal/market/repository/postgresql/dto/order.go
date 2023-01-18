package dto

import (
	"errors"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const (
	NEW = iota
	PROCESSING
	INVALID
	PROCESSED
)

type OrderDTO struct {
	ID		   int64            `db:"order_id,omitempty"`
	Number     string           `db:"order_number,omitempty"`
	Status     int64            `db:"status,omitempty"`
	Accrual    float64          `db:"accrual,omitempty"`
	UploadedAt pgtype.Timestamp `db:"uploaded_at,omitempty"`
	UserID 	   uuid.UUID 		`db:"fk_user_id,omitempty"`
}

func OrdersToModel(ordersDTO []OrderDTO) ([]model.Order, error) {
	var orders []model.Order

	for _, o := range ordersDTO {
		var status string

		switch o.Status {
		case NEW:
			status = "NEW"
		case PROCESSING:
			status = "PROCESSING"
		case INVALID:
			status = "INVALID"
		case PROCESSED:
			status = "PROCESSED"
		default:
			return nil, errors.New("invalid status from DB")
		}

		val, err := o.UploadedAt.Value()
		if err != nil {
			return nil, err
		}

		timeStamp, ok := val.(time.Time)
		if !ok {
			return nil, errors.New("error type assertion timestamp")
		}

		order := model.Order{
			Number:     o.Number,
			Status:     status,
			Accrual:    o.Accrual,
			UploadedAt: timeStamp.Format(time.RFC3339),
		}

		orders = append(orders, order)

	}

	return orders, nil
}
