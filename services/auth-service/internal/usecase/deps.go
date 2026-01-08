package usecase

import "context"

// Интерфейсы для зависимостей, которые реализует Infrastructure
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenManager interface {
	GeneratePair(userID string) (access, refresh string, err error)
}

type UserRegister interface {
	Execute(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
}

type UserLogin interface {
	Execute(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
}

type TokenRefresher interface {
	Execute(ctx context.Context, req *RefreshRequest) (*AuthResponse, error)
}

type PasswordReseter interface {
	Execute(ctx context.Context, req *ResetPasswordRequest) (*AuthResponse, error)
}

type UserLogout interface {
	Execute(ctx context.Context, req *RefreshRequest) error
}
