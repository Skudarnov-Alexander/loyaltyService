package postgresql

import (
	"context"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/auth/repository/postgresql/dto"
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

func createBalance(ctx context.Context, db *sqlx.DB, userID string) error {
	quary := `INSERT INTO balances(fk_user_id)
	VALUES
	($1);`

	tx, err := db.Beginx()
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

func createNewUser(ctx context.Context, db *sqlx.DB, u model.User) error {
	quary := `INSERT INTO users
	VALUES
	($1, $2, $3);`

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.ID, u.Username, u.HashedPass)
	if err != nil {
                if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
                        return auth.ErrUserIsExist
                }
		return err
	}
	
	return tx.Commit()
}

func (p PostrgeSQL) CreateUser(ctx context.Context, u model.User) error {
	if err := createNewUser(ctx, p.db, u); err != nil {
		return err
	}

	if err := createBalance(ctx, p.db, u.ID); err != nil {
		return err
	}

	return nil

}

func (p PostrgeSQL) GetUser(ctx context.Context, username string) (model.User, error) {
	log.Printf("SQL Get user start")
	log.Printf("username: %s", username)

	quary := `SELECT user_id, username, password 
	FROM users 
	WHERE username = $1`

	tx, err := p.db.Beginx()
	if err != nil {
		return model.User{}, err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		log.Printf("PreparexContext err: %s", err.Error())
		return model.User{}, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, username)
	if err != nil {
		log.Printf("QueryxContext err: %s", err.Error())
		return model.User{}, err
	}

	if rows.Err() != nil {
		return model.User{}, rows.Err()
	}

	var user dto.User

	rows.Next()

	err = rows.StructScan(&user)
	if err != nil {
		log.Printf("StructScan err: %s", err.Error())
		return model.User{}, err
	}

	u, err := dto.UserToModel(user)
	if err != nil {
		log.Printf("toModel err: %s", err.Error())
		return model.User{}, err
	}

	return u, tx.Commit()
}
