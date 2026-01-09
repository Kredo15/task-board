package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
	"github.com/Kredo15/task-board/services/auth-service/pkg/grpcutil"
	"github.com/Kredo15/task-board/services/auth-service/pkg/validator"
)

// MapError — вспомогательная функция для перевода доменных ошибок в коды gRPC
func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, domain.ErrInvalidEmail):
		return status.Error(codes.InvalidArgument, "invalid email")
	case errors.Is(err, domain.ErrInvalidPassword):
		return status.Error(codes.Unauthenticated, "invalid password")
	case errors.Is(err, domain.ErrExpiredRefreshToken):
		return status.Error(codes.Unauthenticated, "expired refresh token")
	case errors.Is(err, grpcutil.ErrNoUserID):
		return status.Error(codes.Unauthenticated, "missing user identity")
	case errors.Is(err, domain.ErrEmailAlreadyTaken):
		return status.Error(codes.AlreadyExists, "email already taken")
	case errors.Is(err, domain.ErrInvalidSession):
		return status.Error(codes.Unauthenticated, "invalid session")
	case errors.Is(err, domain.ErrPasswordTooWeak):
		return status.Error(codes.InvalidArgument, "password too weak")
	case errors.Is(err, domain.ErrEmptyHash):
		return status.Error(codes.InvalidArgument, "empty hash")
	case errors.Is(err, validator.ErrInvalidInputData):
		return status.Errorf(codes.InvalidArgument, "invalid input data: %v", err)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
