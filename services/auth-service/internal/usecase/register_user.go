package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
)

type RegisterUseCase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.RefreshTokenRepository
	hasher    PasswordHasher
	tokens    TokenManager
	gen       domain.IDGenerator
	logger    logger.Logger
}

func NewRegisterUseCase(
	ur domain.UserRepository,
	tr domain.RefreshTokenRepository,
	ph PasswordHasher,
	tm TokenManager,
	g domain.IDGenerator,
	log logger.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:  ur,
		tokenRepo: tr,
		hasher:    ph,
		tokens:    tm,
		gen:       g,
		logger:    log,
	}
}

// Register — сценарий регистрации
func (rUC *RegisterUseCase) Execute(ctx context.Context, dto *RegisterRequest) (*RegesterResponse, error) {

	rUC.logger.Info("attempting to register new user")

	email, err := domain.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	exists, err := rUC.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, domain.ErrEmailAlreadyTaken
	}

	// Хешируем пароль
	hash, err := rUC.hasher.Hash(dto.Password)
	if err != nil {
		return nil, err
	}

	// Создаем доменную сущность
	user, err := domain.NewUser(rUC.gen, dto.Email, dto.Username, hash)
	if err != nil {
		return nil, err
	}

	// Сохраняем (через интерфейс репозитория)
	err_create := rUC.userRepo.Create(ctx, user)
	if err_create != nil {
		return nil, fmt.Errorf("failed to create user: %w", err_create)
	}
	// Генерируем токены
	access, refresh, err := rUC.tokens.GeneratePair(user.ID())
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
	if err := rUC.tokenRepo.Save(ctx, tokenEntity); err != nil {
		return nil, err
	}

	return &RegesterResponse{
		UserID:       user.ID(),
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
