package domain

import (
	"time"
)

type User struct {
	ID           UserdID
	Email        Email
	Username     string
	PasswordHash PasswordHash
	CreatedAt    time.Time
}

// NewUser — "фабрика" для создания нового пользователя с базовой валидацией через VO
func NewUser(gen IDGenerator, emailStr, username, passwordHash string) (*User, error) {
	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	hash, err := NewPasswordHash(passwordHash)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           UserdID(gen.Generate()),
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}, nil
}
