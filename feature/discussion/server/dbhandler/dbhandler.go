package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", "postgres://postgres:testing@localhost:5432/postgres")
	if err != nil {
		return nil, fmt.Errorf("connecting to default postgres db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging postgres db: %w", err)
	}

	const dbName = "democracy_routes"

	db.Exec("CREATE DATABASE " + dbName)

	err = db.Close()
	if err != nil {
		return nil, fmt.Errorf("closing postgres db: %w", err)
	}

	db, err = sql.Open("pgx", "postgres://postgres:testing@localhost:5432/"+dbName)
	if err != nil {
		return nil, fmt.Errorf("connecting to %s db: %w", dbName, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging %s db: %w", dbName, err)
	}

	return db, nil
}

func PrepareDB(db *sql.DB) error {
	var aErr error

	if err := createUsersTable(db); err != nil {
		aErr = errors.Join(aErr, err)
	}

	if err := createAuthenticationTable(db); err != nil {
		aErr = errors.Join(aErr, err)
	}

	return aErr
}

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
			id uuid PRIMARY KEY,
			name text UNIQUE NOT NULL
		)`)
	if err != nil {
		return fmt.Errorf("creating users table: %w", err)
	}

	return nil
}

func createAuthenticationTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS authentication (
    		id uuid PRIMARY KEY,
    		hash text NOT NULL
        )`)
	if err != nil {
		return fmt.Errorf("creating authentication table: %w", err)
	}

	return nil
}
