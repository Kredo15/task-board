package domain

import (
	"time"
)

type User struct {
	id           UserdID
	email        Email
	username     string
	passwordHash PasswordHash
	createdAt    time.Time
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
		id:           UserdID(gen.Generate()),
		email:        email,
		username:     username,
		passwordHash: hash,
		createdAt:    time.Now().UTC(),
	}, nil
}

func RestoreUser(id, email, username, password string, createAt time.Time) *User {
	return &User{
		id:           UserdID(id),
		email:        Email(email),
		username:     username,
		passwordHash: PasswordHash(password),
		createdAt:    createAt,
	}
}

func (u *User) ID() string { return string(u.id) }

func (u *User) Email() string { return string(u.email) }

func (u *User) Username() string { return u.username }

func (u *User) Password() string { return string(u.passwordHash) }

func (u *User) CreatedAt() time.Time { return u.createdAt }
