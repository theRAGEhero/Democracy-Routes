package userhandler

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler/model"
)

type Handler struct {
	source map[string]*model.User
}

func New(source map[string]*model.User) (*Handler, error) {
	if source == nil {
		return nil, fmt.Errorf("no source")
	}

	return &Handler{
		source: source,
	}, nil
}

func (h *Handler) Create(params *model.CreateUser) (*model.User, error) {
	var user model.User
	user.ID = uuid.NewString()
	user.Name = params.Name

	h.source[user.ID] = &user

	return &user, nil
}

func (h *Handler) Get(id string) (*model.User, error) {
	user, ok := h.source[id]
	if !ok {
		return nil, semerr.NewNotFoundError(fmt.Errorf("user not found"))
	}

	return user, nil
}
