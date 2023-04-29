package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/jmoiron/sqlx"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return OrderRepository{
		db: db,
	}
}

func (repo *OrderRepository) CheckOrder(ctx context.Context, userID, orderID string) (bool, error) {
	fmt.Println("START check order")
	quary := `SELECT EXISTS 
        (
                SELECT *
                FROM orders
                WHERE order_number = $1 AND fk_user_id = $2
        );`

	tx, err := repo.db.Beginx()
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

func (repo *OrderRepository) InsertOrder(ctx context.Context, userID, orderID string) error {
	fmt.Println("START Insert order")
	quary := `INSERT INTO orders(order_number, fk_user_id) VALUES ($1, $2);`

	tx, err := repo.db.Beginx()
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

func (repo *OrderRepository) SelectOrders(ctx context.Context, userID string) ([]model.Order, error) {
	quary := `SELECT order_number, status, accrual, uploaded_at
	FROM orders
	WHERE fk_user_id = $1
	ORDER BY uploaded_at DESC;`

	tx, err := repo.db.Beginx()
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
