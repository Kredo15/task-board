package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
)

type ResetUseCase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.RefreshTokenRepository
	hasher    PasswordHasher
	tokens    TokenManager
}

func NewResetUseCase(
	ur domain.UserRepository,
	tr domain.RefreshTokenRepository,
	ph PasswordHasher,
	tm TokenManager,
) *ResetUseCase {
	return &ResetUseCase{
		userRepo:  ur,
		tokenRepo: tr,
		hasher:    ph,
		tokens:    tm,
	}
}

func (rUC *ResetUseCase) Execute(ctx context.Context, dto *ResetPasswordRequest) (*AuthResponse, error) {
	// логика смены пароля в БД
	// Ищем пользователя
	userID := domain.NewUserID(dto.ID)
	user, err := rUC.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Проверяем пароль
	if err := rUC.hasher.Compare(user.Password(), dto.OldPassword); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Принудительно разлогиниваем пользователя везде
	err_redis := rUC.tokenRepo.DeleteAllByUserID(ctx, user.ID())
	if err_redis != nil {
		return nil, err_redis
	}
	// Генерируем новую пару токенов
	access, refresh, err := rUC.tokens.GeneratePair(user.ID())
	if err != nil {
		return nil, err
	}
	refreshTTL := time.Hour * 24 * 30 // Заменить на значение из Config
	expiresAt := time.Now().UTC().Add(refreshTTL)

	tokenEntity := &domain.RefreshToken{
		UserID:    user.ID(),
		Token:     refresh, // В идеале здесь тоже хранить хеш
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
	}
	if err := rUC.tokenRepo.Save(ctx, tokenEntity); err != nil {
		return nil, err
	}
	return &AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    expiresAt,
	}, nil
}
