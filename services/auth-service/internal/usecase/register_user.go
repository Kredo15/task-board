package usecase

import (
	"context"
	"fmt"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
)

type RegisterUseCase struct {
	userRepo domain.UserRepository
	hasher   PasswordHasher
	gen      domain.IDGenerator
	logger   logger.Logger
}

func NewRegisterUseCase(
	ur domain.UserRepository,
	ph PasswordHasher,
	g domain.IDGenerator,
	log logger.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: ur,
		hasher:   ph,
		gen:      g,
		logger:   log,
	}
}

// Register — сценарий регистрации
func (rUC *RegisterUseCase) Execute(ctx context.Context, dto *RegisterRequest) error {

	rUC.logger.Info("attempting to register new user")

	email, err := domain.NewEmail(dto.Email)

	exists, err := rUC.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return domain.ErrEmailAlreadyTaken
	}

	// Хешируем пароль
	hash, err := rUC.hasher.Hash(dto.Password)
	if err != nil {
		return err
	}

	// Создаем доменную сущность (бизнес-валидация внутри NewUser)
	user, err := domain.NewUser(rUC.gen, dto.Email, dto.Username, hash)
	if err != nil {
		return err
	}

	// Сохраняем (через интерфейс репозитория)
	return rUC.userRepo.Create(ctx, user)
}
