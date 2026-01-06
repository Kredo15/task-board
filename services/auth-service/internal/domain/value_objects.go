package domain

import (
	"net/mail"
)

type UserdID string

type IDGenerator interface {
	Generate() string
}

type Email string

func NewEmail(v string) (Email, error) {
	if _, err := mail.ParseAddress(v); err != nil {
		return "", ErrInvalidEmail
	}
	return Email(v), nil
}

type PasswordHash string

func NewPasswordHash(hash string) (PasswordHash, error) {
	if len(hash) == 0 {
		return "", ErrEmptyHash
	}
	return PasswordHash(hash), nil
}
