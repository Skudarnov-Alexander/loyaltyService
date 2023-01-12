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
		password text NOT NULL
	);
	
	CREATE TABLE orders
	(
		order_id varchar(32) PRIMARY KEY,
		status smallint DEFAULT 0,
		accrual integer DEFAULT 0,
		uploaded_at timestamp NOT NULL,
		fk_user_id uuid REFERENCES users(user_id) NOT NULL
	);
	
	CREATE TABLE balances
	(
		current_balance integer DEFAULT 0 NOT NULL,
		withdrawn integer,
		fk_user_id uuid REFERENCES users(user_id) NOT NULL UNIQUE
	);
	
	CREATE TABLE purchases
	(
		purchase_id SERIAL PRIMARY KEY,
		sum integer NOT NULL,
		processed_at timestamp NOT NULL,
		fk_user_id uuid REFERENCES users(user_id) NOT NULL
	);
	
	INSERT INTO users
	VALUES
	('db61d134-aa52-49d9-a006-4e82e4d237ca', 'test', 'e1a9b8512ff9c6383e11c59b987aa596b5902b366190849079f0f3b622aaa2d2ef8dbfc3b3e35b266d9237d075e7432cf2afcbac85ec2e649eeb7cb1d5009d64');`

	if _, err := db.Exec(quary); err != nil {
		return err
	}

	log.Print("Tables in DataBase is created")

	return nil
}