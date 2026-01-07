package grpc

// Метод ResetPassword

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
)

// ResetPassword — точка входа для обновления пароля
func (h *Handler) ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error) {
	// Маппинг из Proto в usecase DTO
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("метаданные не найдены")
	}
	values := md.Get("x-user-id")
	if len(values) == 0 {
		return nil, fmt.Errorf("user-id отсутствует в метаданных")
	}

	userID := values[0]
	regDTO := usecase.ResetPasswordRequest{
		UserID:      userID,
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}

	authUC, err := h.deps.ResetPasswordUC.Execute(ctx, &regDTO)
	if err != nil {
		return nil, mapError(err)
	}

	return &authv1.ResetPasswordResponse{
		AccessToken:  authUC.AccessToken,
		RefreshToken: authUC.RefreshToken,
	}, nil
}
