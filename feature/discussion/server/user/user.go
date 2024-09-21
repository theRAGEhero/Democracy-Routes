package user

import "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/user/model"

type Handler struct{}

func (r *Handler) GetUser(id string) (*model.User, error) {
	return &model.User{}, nil
}
