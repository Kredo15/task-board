package grpc

import (
	"errors"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapError — вспомогательная функция для перевода доменных ошибок в коды gRPC
func mapError(err error) error {
	if errors.Is(err, domain.ErrUserNotFound) {
		return status.Error(codes.NotFound, "user not found")
	}
	if errors.Is(err, domain.ErrInvalidEmail) {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	return status.Error(codes.Internal, "internal server error")
}
