package database

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DSN string = "postgres://postgres:postgres@localhost:5432/marketDB"

func New() (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", DSN)
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
		order_id varchar(32) PRIMARY KEY,
		status smallint DEFAULT 0,
		accrual integer DEFAULT 0,
		uploaded_at TIMESTAMP DEFAULT Now(),
		fk_user_id uuid REFERENCES users(user_id) NOT NULL
	);
	
	CREATE TABLE balances
	(
		current_balance real DEFAULT 500 NOT NULL,
		withdrawn real DEFAULT 0,
		fk_user_id uuid REFERENCES users(user_id) NOT NULL,
		UNIQUE(fk_user_id)
	);
	
	CREATE TABLE purchases
	(
		purchase_id SERIAL PRIMARY KEY,
		order_id varchar(32) NOT NULL,
		sum integer NOT NULL,
		processed_at TIMESTAMP DEFAULT Now(),
		fk_user_id uuid REFERENCES users(user_id) NOT NULL,
		UNIQUE(order_id)
	);
	
	INSERT INTO users
	VALUES
	('db61d134-aa52-49d9-a006-4e82e4d237ca', 'test', '083ade633acab7c70de63b24c620eb36b7e388235af30d67568c2f000deb5d7e56d27177c6467d4d1526b425842543e4a3d9136bc014c17220a5f5396c78b3c9');
	
	INSERT INTO orders(order_id, uploaded_at, fk_user_id)
	VALUES
	('657883737','1999-01-08 04:05:06', 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('657887875','1999-01-23 04:05:06', 'db61d134-aa52-49d9-a006-4e82e4d237ca'),
	('657887874','1999-04-23 04:05:06', 'db61d134-aa52-49d9-a006-4e82e4d237ca');`

	if _, err := db.Exec(quary); err != nil {
		return err
	}

	log.Print("Tables in DataBase is created")

	return nil
}