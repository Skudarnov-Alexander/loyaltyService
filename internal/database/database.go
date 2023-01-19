package database

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)


func New(addr string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", addr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Print("postgreSQL is connected")

	return db, nil
}

func InitDB(db *sqlx.DB) error {
	quary := 
	`DROP TABLE IF EXISTS orders;
	DROP TABLE IF EXISTS balances;
	DROP TABLE IF EXISTS purchases;
	DROP TABLE IF EXISTS users;
	
	CREATE TABLE users
	(
		user_id uuid PRIMARY KEY,
		username varchar(16) NOT NULL,
		password text NOT NULL,
		UNIQUE(username)
	);
	
	CREATE TABLE orders
	(	
		order_id SERIAL,
		order_number varchar(32) PRIMARY KEY,
		status smallint DEFAULT 0,
		accrual real DEFAULT NULL,
		uploaded_at TIMESTAMP DEFAULT Now(),
		fk_user_id uuid REFERENCES users(user_id) NOT NULL,
		UNIQUE(order_id)
	);
	
	CREATE TABLE balances
	(
		current_balance real DEFAULT 0 NOT NULL,
		withdrawn real DEFAULT 0,
		fk_user_id uuid REFERENCES users(user_id) NOT NULL,
		UNIQUE(fk_user_id)
	);
	
	CREATE TABLE purchases
	(
		purchase_id SERIAL PRIMARY KEY,
		order_number varchar(32) NOT NULL,
		sum integer NOT NULL,
		processed_at TIMESTAMP DEFAULT Now(),
		fk_user_id uuid REFERENCES users(user_id) NOT NULL,
		UNIQUE(order_number)
	);
	
	INSERT INTO users
	VALUES
	('db61d134-aa52-49d9-a006-4e82e4d237ca', 'test', '083ade633acab7c70de63b24c620eb36b7e388235af30d67568c2f000deb5d7e56d27177c6467d4d1526b425842543e4a3d9136bc014c17220a5f5396c78b3c9');
	
	INSERT INTO orders(order_number, status, fk_user_id)
	VALUES
	('1', 0, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('2', 1, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('3', 0, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('4', 1, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('5', 0, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('6', 0, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('7', 0, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('8', 0, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('9', 1, 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('10', 0, 'db61d134-aa52-49d9-a006-4e82e4d237ca');`

	if _, err := db.Exec(quary); err != nil {
		return err
	}

	log.Print("Tables in DataBase is created")

	return nil
}