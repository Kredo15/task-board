package grpc

// Методы Refresh и Logout

import (
	"context"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
)

// Refresh — точка входа для обновления пары токенов
func (h *Handler) Refresh(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	// Маппинг из Proto в usecase DTO
	RefreshDTO := usecase.RefreshRequest{
		RefreshToken: req.GetRefreshToken(),
	}
	// Вызов Usecase
	authUC, err := h.refreshUC.Execute(ctx, &RefreshDTO)
	if err != nil {
		return nil, mapError(err)
	}

	// Маппинг из DTO в Proto ответ
	return &authv1.RefreshTokenResponse{
		AccessToken:  authUC.AccessToken,
		RefreshToken: authUC.RefreshToken,
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	// Маппинг из Proto в usecase DTO
	RefreshDTO := usecase.RefreshRequest{
		RefreshToken: req.GetToken(),
	}
	// Вызов Usecase
	err := h.logoutUC.Execute(ctx, &RefreshDTO)
	if err != nil {
		return nil, mapError(err)
	}

	// Маппинг из DTO в Proto ответ
	return &authv1.LogoutResponse{}, nil
}
