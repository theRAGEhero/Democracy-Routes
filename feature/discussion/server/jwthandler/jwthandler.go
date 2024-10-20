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

func (h *Handler) Verify(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return h.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("parsing token: %w", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("getting subject: %w", err)
	}

	return subject, nil
}
