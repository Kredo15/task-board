package postgres

import (
	"time"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
)

// userModel — как пользователь выглядит в таблице Postgres
type userModel struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

// Маппер: DB -> Domain
func (m *userModel) toDomain() *domain.User {
	return &domain.User{
		ID:           domain.UserdID(m.ID),
		Email:        domain.Email(m.Email),
		Username:     m.Username,
		PasswordHash: domain.PasswordHash(m.PasswordHash),
		CreatedAt:    m.CreatedAt,
	}
}
