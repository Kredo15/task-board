package domain

import (
	"net/mail"
	"strings"
)

type IDGenerator interface {
	Generate() string
}

type UserdID string

func NewUserID(id string) UserdID {
	return UserdID(id)
}

type Email string

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrInvalidEmail
	}
	return Email(email), nil
}

type PasswordHash string

func NewPasswordHash(hash string) (PasswordHash, error) {
	if len(hash) == 0 {
		return "", ErrEmptyHash
	}
	return PasswordHash(hash), nil
}
