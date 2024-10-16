package jwthandler

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	secret []byte
}

func New(secret []byte) (*Handler, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("no secret")
	}

	return &Handler{
		secret: secret,
	}, nil
}

func (h *Handler) Issue(subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	})

	ss, err := token.SignedString(h.secret)
	if err != nil {
		return "", fmt.Errorf("signing string: %w", err)
	}

	return ss, nil
}
