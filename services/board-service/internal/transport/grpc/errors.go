package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
	"github.com/Kredo15/task-board/services/board-service/pkg/grpcutil"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"
)

func mapErrorToGRPC(err error) error {
	switch {
	case errors.Is(err, board.ErrInvalidBoardTitleEmpty):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, board.ErrInvalidOwnerID):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, validator.ErrInvalidInputData):
		return status.Errorf(codes.InvalidArgument, "invalid input data: %v", err)
	case errors.Is(err, grpcutil.ErrNoUserID):
		return status.Error(codes.Unauthenticated, "missing user identity")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
