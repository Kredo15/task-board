package grpc

import (
	"context"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
)

// Register — точка входа для регистрации
func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// Маппинг из Proto в usecase DTO
	regDTO := usecase.RegisterRequest{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	regUC, err := h.deps.RegisterUC.Execute(ctx, &regDTO)
	if err != nil {
		h.deps.Logger.Error("failed to register user", err)
		return nil, mapError(err)
	}

	return &authv1.RegisterResponse{
		UserId:       regUC.UserID,
		AccessToken:  regUC.AccessToken,
		RefreshToken: regUC.RefreshToken,
	}, nil
}
