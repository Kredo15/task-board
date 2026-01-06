package domain

import (
	"net/mail"
)

type UserdID string

type IDGenerator interface {
	Generate() string
}

// Email — Value Object
type Email struct {
	value string
}

func NewEmail(v string) (Email, error) {
	if _, err := mail.ParseAddress(v); err != nil {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: v}, nil
}

func (e Email) String() string {
	return e.value
}

// PasswordHash — Value Object (хранит уже захешированный пароль)
type PasswordHash struct {
	value string
}

func NewPasswordHash(hash string) (PasswordHash, error) {
	if len(hash) == 0 {
		return PasswordHash{}, ErrEmptyHash
	}
	return PasswordHash{value: hash}, nil
}

func (p PasswordHash) String() string {
	return p.value
}
