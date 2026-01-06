package postgres

import (
	"time"
)

// userModel — как пользователь выглядит в таблице Postgres
type userModel struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
