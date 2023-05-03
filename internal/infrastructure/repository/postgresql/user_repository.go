package postgresql

import (
	"context"
	"fmt"
	"log"

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
		`INSERT INTO users (username, password) 
		VALUES ($1, $2)
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

	rows, err := stmt.QueryxContext(ctx, u.Username, u.HashedPass)
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

	if err := rows.StructScan(&user); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	id, err := user.ID.Value()
	if err != nil {
		return "", err
	}

	return id.(string), nil
}

func (repo *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	quary := `SELECT id, username, password
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
		fmt.Printf("rows ERR: %v\n", rows.Err())
		return model.User{}, rows.Err()
	}

	var user dto.User

	if ok := rows.Next(); !ok {
		fmt.Printf("user is not found: %v\n", rows.Err())
		return model.User{}, error2.ErrUserNotFound
	}

	err = rows.StructScan(&user)
	if err != nil {
		log.Printf("StructScan err: %s", err.Error())
		return model.User{}, err
	}

	u, err := user.ToModel()
	if err != nil {
		log.Printf("toModel err: %s", err.Error())
		return model.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, err
	}

	return u, nil
}
