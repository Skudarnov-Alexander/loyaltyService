package postgresql

import (
	"context"
	"errors"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

const LIMIT_ORDERS int64 = 5


func (p *PostrgeSQL) TakeOrdersForProcess(ctx context.Context) ([]model.Accrual, error) {
	quary := `SELECT order_number
	FROM orders
	WHERE status = 0
	ORDER BY uploaded_at
	LIMIT $1;`

	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var accrualsDTO []dto.Accrual
	err = stmt.Select(&accrualsDTO, LIMIT_ORDERS)
	if err != nil {
		return nil, err
	}

	var accruals []model.Accrual
	for _, a := range accrualsDTO {
		accrual := dto.AccrualToModel(a)
		accruals = append(accruals, accrual)
	}

	return accruals, tx.Commit()
}

func (p *PostrgeSQL) ChangeStatusOrdersForProcess(ctx context.Context, accruals ...model.Accrual) error {
	quary := `UPDATE orders
	SET
			status = 2
	WHERE order_number = $1
	RETURNING order_number, status;`

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, a := range accruals {
		rows, err := stmt.QueryxContext(ctx, a.Number)
		if err != nil {
			tx.Rollback()
			return err
		}

		if rows.Err() != nil {
			return rows.Err()
		}

		rows.Next()

		var accrualsDTO dto.Accrual
		err = rows.StructScan(&accrualsDTO)
		if err != nil {
			return err
		}

		log.Printf("Change status: %+v", accrualsDTO)
		rows.Close()
	}

	return tx.Commit()
}

func (p *PostrgeSQL) UpdateStatusProcessedOrders(ctx context.Context, a model.Accrual) error { //TODO батчи?
        var status int64

	switch a.Status {
	case "INVALID":
		status = 2
	case "PROCESSED":
		status = 3
	default:
		return errors.New("invalid status from DB")
	}

        
	quary := `UPDATE orders
	SET
		status = $1,
                accrual = $2
	WHERE order_number = $3
	RETURNING order_number, status;`

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	
	rows, err := stmt.QueryxContext(ctx, status, a.Accrual, a.Number)
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	rows.Next()

	var accrualsDTO dto.Accrual
	err = rows.StructScan(&accrualsDTO)
	if err != nil {
		return err
	}

	log.Printf("Change status: %+v", accrualsDTO)
	rows.Close()
	

	return tx.Commit()
}
