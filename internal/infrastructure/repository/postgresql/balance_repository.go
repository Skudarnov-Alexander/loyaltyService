package postgresql

import (
	"context"
	"errors"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type BalanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{
		db: db,
	}
}

func (repo *BalanceRepository) Сreate(ctx context.Context, userID string) error {
	quary := `INSERT INTO balances(fk_user_id) VALUES ($1);`

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

	_, err = stmt.ExecContext(ctx, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *BalanceRepository) GetByUserId(ctx context.Context, id string) (model.Balance, error) {
	quary := `SELECT current_balance, withdrawn
	FROM balances
	WHERE fk_user_id = $1;`

	tx, err := repo.db.Beginx()
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

	err = stmt.Get(&balanceDTO, id)
	if err != nil {
		return model.Balance{}, err
	}

	balance := dto.BalanceToModel(balanceDTO)

	return balance, tx.Commit()
}

func Update(ctx context.Context, db *sqlx.DB, userID string, b dto.BalanceDTO) (dto.BalanceDTO, error) {
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
