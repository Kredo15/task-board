package usecase

import (
	"context"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
	"github.com/Kredo15/task-board/services/auth-service/pkg/logger"
)

type LogoutUseCase struct {
	tokenRepo domain.TokenRepository
	logger    logger.Logger
}

func NewLogoutUseCase(
	tr domain.TokenRepository,
	l logger.Logger,
) *LogoutUseCase {
	return &LogoutUseCase{
		tokenRepo: tr,
		logger:    l,
	}
}

func (lUC *LogoutUseCase) Execute(ctx context.Context, dto *RefreshRequest) error {
	return lUC.tokenRepo.DeleteByToken(ctx, dto.RefreshToken)
}
