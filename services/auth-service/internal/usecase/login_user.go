package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
)

type LoginUseCase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.RefreshTokenRepository
	hasher    PasswordHasher
	tokens    TokenManager
	logger    logger.Logger
}

func NewLoginUseCase(
	ur domain.UserRepository,
	tr domain.RefreshTokenRepository,
	ph PasswordHasher,
	tm TokenManager,
	log logger.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:  ur,
		tokenRepo: tr,
		hasher:    ph,
		tokens:    tm,
		logger:    log,
	}
}

// Login — сценарий входа
func (lUC *LoginUseCase) Execute(ctx context.Context, dto *LoginRequest) (*AuthResponse, error) {
	email, err := domain.NewEmail(dto.Email)
	if err != nil {
		return nil, domain.ErrInvalidEmail
	}
	// Ищем пользователя
	user, err := lUC.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	// Проверяем пароль
	if err := lUC.hasher.Compare(user.Password(), dto.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Генерируем токены
	access, refresh, err := lUC.tokens.GeneratePair(user.ID())
	if err != nil {
		return nil, err
	}

	refreshTTL := time.Hour * 24 * 30 // Заменить на значение из Config
	expiresAt := time.Now().UTC().Add(refreshTTL)

	// Сохраняем refresh токен для возможности отзыва
	tokenEntity := &domain.RefreshToken{
		UserID:    user.ID(),
		Token:     refresh, // В идеале здесь тоже хранить хеш
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
	}
	if err := lUC.tokenRepo.Save(ctx, tokenEntity); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    expiresAt,
	}, nil
}
