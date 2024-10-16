package authenticationhandler

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) (*Handler, error) {
	if db == nil {
		return nil, fmt.Errorf("no db")
	}

	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) SetPassword(id string, password string) error {
	phash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generating password hash: %w", err)
	}

	_, err = h.db.Exec("INSERT INTO authentication (id, hash) VALUES ($1, $2) "+
		"ON CONFLICT (id) DO UPDATE SET hash = EXCLUDED.hash", id, phash)
	if err != nil {
		return fmt.Errorf("setting password: %w", err)
	}

	return nil
}

func (h *Handler) Authenticate(id string, pass string) bool {
	var phash string
	if err := h.db.QueryRow("SELECT hash FROM authentication WHERE id = $1", id).Scan(&phash); err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(phash), []byte(pass)) == nil
}
