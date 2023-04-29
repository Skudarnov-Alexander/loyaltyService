package postgresql

import (
	"context"

	error2 "github.com/Skudarnov-Alexander/loyaltyService/internal/error"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/infrastructure/repository/postgresql/dto"
	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u model.User) (string, error) {
	quary :=
		`INSERT INTO users VALUES ($1, $2, $3)
		RETURNING *;`

	tx, err := repo.db.Beginx()
	if err != nil {
		return "", err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, u.ID, u.Username, u.HashedPass)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			return "", error2.ErrUserIsExist
		}
		return "", err
	}

	if rows.Err() != nil {
		return "", rows.Err()
	}

	rows.Next()

	var user dto.User

	if err := rows.Scan(&user); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return user.ID.String(), nil
}

/*
func (p PostrgeSQL) CreateUser(ctx context.Context, u model.User) error {
	if err := createNewUser(ctx, p.db, u); err != nil {
		return err
	}

	if err := createBalance(ctx, p.db, u.ID); err != nil {
		return err
	}

	return nil

}
*/

/*
func (repo *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	log.Printf("SQL Get user start")
	log.Printf("username: %s", username)

	quary := `SELECT user_id, username, password
	FROM users
	WHERE username = $1`

	tx, err := repo.db.Beginx()
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
*/
