package usecase

import (
	"context"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
)

type LogoutUseCase struct {
	tokenRepo domain.RefreshTokenRepository
}

func NewLogoutUseCase(
	tr domain.RefreshTokenRepository,
) *LogoutUseCase {
	return &LogoutUseCase{
		tokenRepo: tr,
	}
}

func (lUC *LogoutUseCase) Execute(ctx context.Context, dto *RefreshRequest) error {
	return lUC.tokenRepo.DeleteByToken(ctx, dto.RefreshToken)
}
