package dto

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	NEW int64 = iota
	PROCESSING
	INVALID
	PROCESSED
)

type Order struct {
	ID         int64            `db:"id,omitempty"`
	Number     string           `db:"order_number,omitempty"`
	Status     int64            `db:"status,omitempty"`
	Accrual    sql.NullFloat64  `db:"accrual,omitempty"`
	UploadedAt pgtype.Timestamp `db:"uploaded_at,omitempty"`
	UserID     uuid.UUID        `db:"fk_user_id,omitempty"`
}

func OrderToModel(ordersDTO ...Order) ([]model.Order, error) {
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
			return nil, fmt.Errorf("OrderTOModel invalid status from DB %+v", 0)
		}

		val, err := o.UploadedAt.Value()
		if err != nil {
			return nil, err
		}

		timeStamp, ok := val.(time.Time)
		if !ok {
			return nil, errors.New("error type assertion timestamp")
		}

		order := model.Order{ //TODO pointer in model
			Number:     o.Number,
			Status:     status,
			Accrual:    o.Accrual.Float64,
			UploadedAt: timeStamp.Format(time.RFC3339),
		}

		orders = append(orders, order)

	}

	return orders, nil
}
