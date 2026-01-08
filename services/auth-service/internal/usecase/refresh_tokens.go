package usecase

import (
	"context"
	"time"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
)

type RefreshUseCase struct {
	tokenRepo domain.TokenRepository
	tokens    TokenManager
	logger    logger.Logger
}

func NewRefreshUseCase(
	tr domain.TokenRepository,
	tm TokenManager,
	l logger.Logger,
) *RefreshUseCase {
	return &RefreshUseCase{
		tokenRepo: tr,
		tokens:    tm,
		logger:    l,
	}
}

func (rUC *RefreshUseCase) Execute(ctx context.Context, dto *RefreshRequest) (*AuthResponse, error) {
	// Проверяем наличие токена в Redis
	token, err := rUC.tokenRepo.GetByToken(ctx, dto.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Если токен валиден, удаляем старый
	_ = rUC.tokenRepo.DeleteByToken(ctx, dto.RefreshToken)

	// Генерируем новую пару токенов
	access, refresh, err := rUC.tokens.GeneratePair(token.UserID)
	if err != nil {
		return nil, err
	}

	refreshTTL := time.Hour * 24 * 30 // Заменить на значение из Config
	expiresAt := time.Now().UTC().Add(refreshTTL)

	tokenEntity := &domain.RefreshToken{
		UserID:    token.UserID,
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
