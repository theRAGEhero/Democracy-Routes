package authandler

import (
	"database/sql"
	"fmt"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) (*Handler, error) {
	if db == nil {
		return nil, semerr.NewInternalServerError(fmt.Errorf("no db"))
	}

	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) SetPassword(id string, password string) error {
	phash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return semerr.NewBadRequestError(fmt.Errorf("generating password hash: %w", err))
	}

	_, err = h.db.Exec("INSERT INTO auth (id, hash) VALUES ($1, $2) " +
		"ON CONFLICT (id) DO UPDATE SET hash = EXCLUDED.hash", id, phash)
	if err != nil {
		return semerr.NewInternalServerError(fmt.Errorf("setting password: %w", err))
	}

	return nil
}
