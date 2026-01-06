package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email Email) (*User, error)
	GetByID(ctx context.Context, id UserdID) (*User, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
	UpdatePassword(ctx context.Context, passHash PasswordHash) error
}

type RefreshTokenRepository interface {
	// Save сохраняет новый токен или обновляет существующий
	Save(ctx context.Context, token *RefreshToken) error

	// GetByToken ищет токен для проверки при обновлении access-пары
	GetByToken(ctx context.Context, token string) (*RefreshToken, error)

	// DeleteByUserID используется для выхода со всех устройств (Logout)
	DeleteAllByUserID(ctx context.Context, userID string) error

	// Delete используется для отзыва конкретного токена
	DeleteByToken(ctx context.Context, token string) error
}
