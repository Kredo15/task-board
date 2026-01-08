package grpc

import (
	"context"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
)

// Login — точка входа для аутентификации
func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	// Маппинг из Proto в usecase DTO
	loginDTO := usecase.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	// Вызов Usecase
	authUC, err := h.loginUC.Execute(ctx, &loginDTO)
	if err != nil {
		return nil, mapError(err)
	}

	// Маппинг из DTO в Proto ответ
	return &authv1.LoginResponse{
		AccessToken:  authUC.AccessToken,
		RefreshToken: authUC.RefreshToken,
	}, nil
}
