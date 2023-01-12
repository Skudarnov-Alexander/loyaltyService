package postgresql

import (
	"context"
	"errors"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"
	//"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
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

// TODO close db conn

func checkUniqueUser(ctx context.Context, db *sqlx.DB, username string) (bool, error) {
	var ok bool
	quary := `SELECT EXISTS
	(
		SELECT username
		FROM users
		WHERE username = $1
	);`

	tx, err := db.Beginx()
	if err != nil {
		return ok, err
	}

	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, quary)
	if err != nil {
		return ok, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, username)
	if err != nil {
		return ok, err
	}

	rows.Next()

	err = rows.Scan(&ok)
	if err != nil {
		return ok, err
	}
	log.Printf("OK %v", ok)

	return ok, tx.Commit()
}

func (p PostrgeSQL) CreateUser(ctx context.Context, u model.User) error {
	ok, err := checkUniqueUser(ctx, p.db, u.Username)
	if err != nil {
		return err
	}

	if ok {
		return errors.New("user is exist")
	}

	quary := `INSERT INTO users
	VALUES
	($1, $2, $3);`

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

	if _, err := stmt.QueryContext(ctx, u.ID, u.Username, u.Password); err != nil {
		return err
	}
	
	return tx.Commit()

}

type User struct {
	UserID uuid.UUID `db:"user_id"`
	Username string  `db:"username"`
	Password string	 `db:"password"`
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

	row, err := stmt.QueryxContext(ctx, username)
	if err != nil {
		log.Printf("QueryxContext err: %s", err.Error())
		return model.User{}, err
	}
	
	var user User
	
	row.Next()
	err = row.StructScan(&user)
	if err != nil {
		log.Printf("StructScan err: %s", err.Error())
		return model.User{}, err
	}
	
	u, err := toModel(user)
	if err != nil {
		log.Printf("toModel err: %s", err.Error())
		return model.User{}, err
	}

	return u, tx.Commit()
}

func toModel (u User) (model.User, error) {
	uuid, err := u.UserID.Value()
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:       uuid.(string),
		Username: u.Username,
		Password: u.Password,
	}, nil
}

/*
func toModel (u User) (model.User, error) {


	return model.User{
		ID:       u.UserID,
		Username: u.Username,
		Password: u.Password,
	}, nil
}
*/