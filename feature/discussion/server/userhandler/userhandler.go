package userhandler

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler/model"
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

func (h *Handler) Create(params *model.CreateUser) (*model.User, error) {
	var user model.User
	user.ID = uuid.NewString()
	user.Name = params.Name

	_, err := h.db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", user.ID, user.Name)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return &user, nil
}

func (h *Handler) Get(id string) (*model.User, error) {
	var user model.User

	err := h.db.
		QueryRow("SELECT id, name FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no such user: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	return &user, nil
}

func (h *Handler) GetByName(name string) (*model.User, error) {
	var user model.User

	err := h.db.
		QueryRow("SELECT id, name FROM users WHERE name = $1", name).
		Scan(&user.ID, &user.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no such user: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	return &user, nil
}
