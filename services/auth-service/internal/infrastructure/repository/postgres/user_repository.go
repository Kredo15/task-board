package postgres

import (
	"context"
	"errors"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, username, password_hash, created_at) 
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query,
		user.ID(),
		user.Email(),
		user.Username(),
		user.Password(),
		user.CreatedAt(),
	)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at 
		FROM users 
		WHERE email = $1
	`

	var m userModel
	err := r.db.QueryRow(ctx, query, email).Scan(
		&m.ID, &m.Email, &m.Username, &m.PasswordHash, &m.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return domain.RestoreUser(m.ID, m.Email, m.Username, m.PasswordHash, m.CreatedAt), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id domain.UserdID) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, created_at 
		FROM users 
		WHERE id = $1
	`

	var m userModel
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Email, &m.Username, &m.PasswordHash, &m.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound // Используем ошибку pgx
		}
		return nil, err
	}

	return domain.RestoreUser(m.ID, m.Email, m.Username, m.PasswordHash, m.CreatedAt), nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
