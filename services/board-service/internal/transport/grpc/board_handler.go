package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
	usecase "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"
)

type BoardHandler struct {
	boardUsecase usecase.CreateBoardUseCase
	boardv1.UnimplementedBoardServiceServer
	validate *validator.Validator
	log      logger.Logger
}

func NewBoardHandler(bUC usecase.CreateBoardUseCase) *BoardHandler {
	return &BoardHandler{
		boardUsecase: bUC,
	}
}

func (h *BoardHandler) CreateBoard(
	ctx context.Context,
	req *boardv1.CreateBoardRequest,
) (*boardv1.CreateBoardResponse, error) {

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}
	h.log.Info("creating board user %s", userID)

	boardDTO := usecase.CreateBoardRequest{
		Title:       req.Title,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := h.validate.Struct(boardDTO); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err)
	}

	createdBoard, err := h.boardUsecase.Execute(ctx, &boardDTO)
	if err != nil {
		h.log.Error("CreateBoard usecase error", err)
		return nil, convertDomainErrorToGRPC(err)
	}

	h.log.Info("board retrieved successfully %s", createdBoard.ID)
	return &boardv1.CreateBoardResponse{
		Board: mapDomainBoardToProto(createdBoard),
	}, nil
}

func convertDomainErrorToGRPC(err error) error {
	switch {
	case errors.Is(err, board.ErrInvalidBoardTitleEmpty):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, board.ErrInvalidOwnerID):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
