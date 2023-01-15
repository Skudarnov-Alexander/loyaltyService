package dto

import (
	"errors"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/jackc/pgtype"
)

const (
	NEW = iota
	PROCESSING
	INVALID
	PROCESSED
)

type OrderDTO struct {
	Number     string           `db:"order_id"`
	Status     int64            `db:"status"`
	Accrual    int64            `db:"accrual"`
	UploadedAt pgtype.Timestamp `db:"uploaded_at"`
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
