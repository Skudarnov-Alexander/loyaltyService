package postgresql

import (
	"context"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/market/repository/postgresql/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type WithdrawnRepository struct {
	db *sqlx.DB
}

func NewWithdrawnRepository(db *sqlx.DB) WithdrawnRepository {
	return WithdrawnRepository{
		db: db,
	}
}

func (repo *WithdrawnRepository) SelectAllByUserId(ctx context.Context, id string) ([]model.Withdrawn, error) {
	quary := `SELECT order_number, sum, processed_at
	FROM purchases
	WHERE fk_user_id = $1
	ORDER BY processed_at;`

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

	var wDTO []dto.WithdrawnDTO

	err = stmt.Select(&wDTO, id)
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
/*
func (repo *WithdrawnRepository) ProcessWithdrawn(ctx context.Context, userID string, w model.Withdrawn) (model.Balance, error) {
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
*/

func (repo *WithdrawnRepository) insertWithdrawal(ctx context.Context, id string, w dto.WithdrawnDTO) error {
	quary := `INSERT INTO purchases(order_number, sum, fk_user_id)
	VALUES
	($1, $2, $3);`

	tx, err := repo.db.Begin()
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
	if _, err := stmt.ExecContext(ctx, w.OrderID, w.Sum, id); err != nil {
		return err
	}

	return tx.Commit()
}
