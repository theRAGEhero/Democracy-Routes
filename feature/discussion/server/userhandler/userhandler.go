package userhandler

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler/model"
)

type Handler struct {
	db     *sql.DB
	source map[string]*model.User
}

func New(db *sql.DB, source map[string]*model.User) (*Handler, error) {
	if db == nil {
		return nil, semerr.NewInternalServerError(fmt.Errorf("no db"))
	}

	return &Handler{
		db:     db,
		source: source,
	}, nil
}

func (h *Handler) Create(params *model.CreateUser) (*model.User, error) {
	var user model.User
	user.ID = uuid.NewString()
	user.Name = params.Name

	if h.db == nil {
		h.source[user.ID] = &user
	} else {
		_, err := h.db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", user.ID, user.Name)
		if err != nil {
			return nil, semerr.NewInternalServerError(fmt.Errorf("creating user: %w", err))
		}
	}

	return &user, nil
}

func (h *Handler) Get(id string) (*model.User, error) {
	var user model.User

	if h.db == nil {
		user, ok := h.source[id]
		if !ok {
			return nil, semerr.NewNotFoundError(fmt.Errorf("user not found"))
		}

		return user, nil
	} else {
		err := h.db.
			QueryRow("SELECT id, name FROM users WHERE id = $1", id).
			Scan(&user.ID, &user.Name)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, semerr.NewNotFoundError(fmt.Errorf("no such user: %w", err))
		} else if err != nil {
			return nil, semerr.NewInternalServerError(fmt.Errorf("getting user: %w", err))
		}
	}

	return &user, nil
}
