package domain

import "errors"

var (
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrEmptyHash           = errors.New("hash cannot be empty")
	ErrPasswordTooWeak     = errors.New("password is too weak")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyTaken   = errors.New("the email address is already taken")
	ErrInvalidSession      = errors.New("invalid session")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrExpiredRefreshToken = errors.New("expired refresh token")
)
