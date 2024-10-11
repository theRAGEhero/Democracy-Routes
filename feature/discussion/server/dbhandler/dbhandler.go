package dbhandler

import (
	"database/sql"
	"errors"
	"fmt"
)

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
