package postgresql

import (
	"context"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DSN string = "postgres://postgres:postgres@localhost:5432/marketDB"

type PostrgeSQL struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*PostrgeSQL, error) {
	return &PostrgeSQL{
		db: db,
	}, nil
}

func (p *PostrgeSQL) InsertOrder(ctx context.Context, userID string, order model.Order) error {
	quary := `INSERT INTO orders(order_id, uploaded_at, fk_user_id)
	VALUES
	($1, $2, $3);`

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, order.Number, order.UploadedAt, userID); err != nil {
		return err
	}
	
	return tx.Commit()
}

func (p *PostrgeSQL) SelectOrdersByUser(ctx context.Context, userID string, order model.Order) error {
	quary := `INSERT INTO orders(order_id, uploaded_at, fk_user_id)
	VALUES
	($1, $2, $3);`

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, order.Number, order.UploadedAt, userID); err != nil {
		return err
	}
	
	return tx.Commit()
}

func (p *PostrgeSQL) SelectBalanceByUser(ctx context.Context, userID string, order model.Order) error {
	quary := `INSERT INTO orders(order_id, uploaded_at, fk_user_id)
	VALUES
	($1, $2, $3);`

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, order.Number, order.UploadedAt, userID); err != nil {
		return err
	}
	
	return tx.Commit()
}

func (p *PostrgeSQL) SelectWithdrawalsByUser(ctx context.Context, userID string, order model.Order) error {
	quary := `INSERT INTO orders(order_id, uploaded_at, fk_user_id)
	VALUES
	($1, $2, $3);`

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, order.Number, order.UploadedAt, userID); err != nil {
		return err
	}
	
	return tx.Commit()
}

func (p *PostrgeSQL) InsertWithdrawal(ctx context.Context, userID string, order model.Order) error {
	quary := `INSERT INTO orders(order_id, uploaded_at, fk_user_id)
	VALUES
	($1, $2, $3);`

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, order.Number, order.UploadedAt, userID); err != nil {
		return err
	}
	
	return tx.Commit()
}




