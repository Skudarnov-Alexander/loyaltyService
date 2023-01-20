package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
)

const limitOrdersBatch int64 = 20

func (p *PostrgeSQL) TakeOrdersForProcess(ctx context.Context) ([]model.Accrual, error) {
	quary := `SELECT order_number, fk_user_id
	FROM orders
	WHERE status = 0 OR status = 2
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
	err = stmt.Select(&accrualsDTO, limitOrdersBatch)
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
			status = 1
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

		log.Printf("Change status for processing: %+v", accrualsDTO)
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

	log.Printf("Change order after processing: %+v", accrualsDTO)
	rows.Close()

	return tx.Commit()
}

func (p *PostrgeSQL) UpdateBalanceProcessedOrders(ctx context.Context, a model.Accrual) error {
        log.Printf("Зашел в изменение баланса %+v", a)
	quary := `UPDATE balances
        SET
	        current_balance = current_balance + $1
        WHERE fk_user_id = $2
        RETURNING current_balance;`

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, a.Accrual, a.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows.Err() != nil {
                fmt.Print("зашел в rows.Err")
		return rows.Err()
	}

	rows.Next()

	var balanceDTO dto.BalanceDTO
	err = rows.StructScan(&balanceDTO)
	if err != nil {
		return err
	}

	log.Printf("Change balance. Add: %f Current: %f", a.Accrual, balanceDTO.Current)
	rows.Close()

	return tx.Commit()
}
