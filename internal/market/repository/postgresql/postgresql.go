package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostrgeSQL struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) PostrgeSQL {
	return PostrgeSQL{
		db: db,
	}
}

func (p PostrgeSQL) CheckOrder(ctx context.Context, userID, orderID string) (bool, error) {
        fmt.Println("START check order")
	quary := `SELECT EXISTS 
        (
                SELECT *
                FROM orders
                WHERE order_number = $1 AND fk_user_id = $2
        );`

	tx, err := p.db.Beginx()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var ok bool
	if err := stmt.Get(&ok, orderID, userID); err != nil {
		return false, err
	}

	return ok, tx.Commit()

}

func (p PostrgeSQL) InsertOrder(ctx context.Context, userID, orderID string) error {
        fmt.Println("START Insert order")
	quary := `INSERT INTO orders(order_number, fk_user_id)
	VALUES
	($1, $2);`

	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, orderID, userID); err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			fmt.Println("OK")
			return market.ErrOrderIsExist
		}
		fmt.Println("not OK")
		return err
	}

        fmt.Println("not not OK")

	return tx.Commit()
}

func (p PostrgeSQL) SelectOrders(ctx context.Context, userID string) ([]model.Order, error) {
	quary := `SELECT order_number, status, accrual, uploaded_at
	FROM orders
	WHERE fk_user_id = $1
	ORDER BY uploaded_at DESC;`

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

	var ordersDTO []dto.Order

	err = stmt.Select(&ordersDTO, userID)
	if err != nil {
		return nil, err
	}

	orders, err := dto.OrderToModel(ordersDTO...)
	if err != nil {
		return nil, err
	}

	log.Printf("OrdersDTO: %+v", ordersDTO)
	log.Printf("Orders: %+v", orders)

	return orders, tx.Commit()
}

func (p PostrgeSQL) SelectBalance(ctx context.Context, userID string) (model.Balance, error) {
	quary := `SELECT current_balance, withdrawn
	FROM balances
	WHERE fk_user_id = $1;`

	tx, err := p.db.Beginx()
	if err != nil {
		return model.Balance{}, err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return model.Balance{}, err
	}

	defer stmt.Close()

	var balanceDTO dto.BalanceDTO

	err = stmt.Get(&balanceDTO, userID)
	if err != nil {
		return model.Balance{}, err
	}

	balance := dto.BalanceToModel(balanceDTO)

	return balance, tx.Commit()
}

func insertWithdrawal(ctx context.Context, db *sqlx.DB, userID string, w dto.WithdrawnDTO) error {
	quary := `INSERT INTO purchases(order_number, sum, fk_user_id)
	VALUES
	($1, $2, $3);`

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	log.Printf("wDTO INSERT %+v", w)
	if _, err := stmt.ExecContext(ctx, w.OrderID, w.Sum, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func loadNewBalance(ctx context.Context, db *sqlx.DB, userID string, w dto.WithdrawnDTO) (dto.BalanceDTO, error) {
	quary := `SELECT current_balance, withdrawn
	FROM balances
	WHERE fk_user_id = $1;`

	var b dto.BalanceDTO

	err := db.Get(&b, quary, userID)
	if err != nil {
		return dto.BalanceDTO{}, err
	}

	if b.Current < w.Sum {
		return dto.BalanceDTO{}, errors.New("user has not enogh balance")
	}

	b.Current -= w.Sum
	b.Withdrawn += w.Sum

	log.Printf("newBalance %v", b)

	return b, nil
}

func updateBalance(ctx context.Context, db *sqlx.DB, userID string, b dto.BalanceDTO) (dto.BalanceDTO, error) {
	quary := `UPDATE balances
        SET
                current_balance = $1,
                withdrawn = $2
        WHERE fk_user_id = $3
        RETURNING current_balance, withdrawn;`

	tx, err := db.Beginx()
	if err != nil {
		return dto.BalanceDTO{}, err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return dto.BalanceDTO{}, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, b.Current, b.Withdrawn, userID)
	if err != nil {
		return dto.BalanceDTO{}, err
	}

	rows.Next()

	var bNew dto.BalanceDTO
	if err := rows.StructScan(&bNew); err != nil {
		log.Printf("structScat error %s", err.Error())
		return dto.BalanceDTO{}, err
	}

	if rows.Err() != nil {
		return dto.BalanceDTO{}, rows.Err()
	}

	return bNew, tx.Commit()
}

func (p PostrgeSQL) ProcessWithdrawn(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error) {
	wDTO := dto.WithdrawnToDTO(w)

	bDTO, err := loadNewBalance(ctx, p.db, userID, wDTO)
	if err != nil {
		return model.Balance{}, err
	}

	if err := insertWithdrawal(ctx, p.db, userID, wDTO); err != nil {
		return model.Balance{}, err
	}

	bNew, err := updateBalance(ctx, p.db, userID, bDTO)
	if err != nil {
		return model.Balance{}, err
	}

	b := dto.BalanceToModel(bNew)

	return b, nil
}

func (p PostrgeSQL) SelectWithdrawals(ctx context.Context, userID string) ([]model.Withdrawn, error) {
	quary := `SELECT order_number, sum, processed_at
	FROM purchases
	WHERE fk_user_id = $1
	ORDER BY processed_at;`

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

	var wDTO []dto.WithdrawnDTO

	err = stmt.Select(&wDTO, userID)
	if err != nil {
		return nil, err
	}

	withdrawns, err := dto.WithdrawnsToModel(wDTO...)
	if err != nil {
		return nil, err
	}

	log.Printf("Withdrawns: %v", withdrawns)

	return withdrawns, tx.Commit()
}
