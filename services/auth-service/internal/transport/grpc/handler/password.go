package grpc

// Метод ResetPassword

import (
	"context"

	authv1 "github.com/Kredo15/task-board/protos/gen/go/auth/v1"
	"github.com/Kredo15/task-board/services/auth-service/internal/usecase"
	"github.com/Kredo15/task-board/services/auth-service/pkg/grpcutil"
)

// ResetPassword — точка входа для обновления пароля
func (h *Handler) ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error) {

	userID, err := grpcutil.GetUserID(ctx)
	if err != nil {
		h.log.Error("missing user identity", err)
		return nil, mapError(err)
	}
	// Маппинг из Proto в usecase DTO
	regDTO := usecase.ResetPasswordRequest{
		UserID:      userID,
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}
	if err := h.valid.Struct(regDTO); err != nil {
		h.log.Error("invalid input data", err)
		return nil, mapError(err)
	}
	authUC, err := h.resetUC.Execute(ctx, &regDTO)
	if err != nil {
		h.log.Error("failed to reset password", err)
		return nil, mapError(err)
	}

	return &authv1.ResetPasswordResponse{
		AccessToken:  authUC.AccessToken,
		RefreshToken: authUC.RefreshToken,
	}, nil
}
