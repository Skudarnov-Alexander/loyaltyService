package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Skudarnov-Alexander/loyaltyService/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	_, err = stmt.QueryxContext(ctx, userID)
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

	_, err = stmt.QueryxContext(ctx, u.ID, u.Username, u.Password)
	log.Printf("%T", err)
	log.Printf("err: %v\n", err)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		fmt.Println(pgErr.Message) // => syntax error at end of input
		fmt.Println(pgErr.Code)    // => 42601
	}

	if err != nil {
		return err
	}
	/*
			if err != nil {
		                //pgerrcode.IsIntegrityConstraintViolation(err.Code)
				if err, ok := err.(pgx.PgError); ok && err.Code == pgerrcode.UniqueViolation {
		                        log.Print("pgerr")
					return errors.New("user is exist")
				}
		                log.Print("not pgerr")

				return err
			  }
	*/
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

type User struct {
	UserID   uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
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

func toModel(u User) (model.User, error) {
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

// не используется

/*
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
*/
