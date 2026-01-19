package grpc

import (
	"context"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	usecase "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
	"github.com/Kredo15/task-board/services/board-service/pkg/grpcutil"
)

type CreateBoard interface {
	Execute(ctx context.Context, cmd *usecase.CreateBoardRequest) (*usecase.BoardResponse, error)
}

func (h *BoardHandler) Create(
	ctx context.Context,
	req *boardv1.CreateBoardRequest,
) (*boardv1.CreateBoardResponse, error) {

	userID, err := grpcutil.GetUserID(ctx)
	if err != nil {
		h.log.Error("missing user identity", err)
		return nil, mapErrorToGRPC(err)
	}

	h.log.Info("creating board user %s", userID)

	boardDTO := usecase.CreateBoardRequest{
		Title:       req.Title,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := h.validate.Struct(boardDTO); err != nil {
		h.log.Error("invalid input", err)
		return nil, mapErrorToGRPC(err)
	}

	createdBoard, err := h.createBoardUC.Execute(ctx, &boardDTO)
	if err != nil {
		h.log.Error("createBoard usecase error", err)
		return nil, mapErrorToGRPC(err)
	}

	h.log.Info("board retrieved successfully %s", createdBoard.ID)
	return &boardv1.CreateBoardResponse{
		Board: mapDTOBoardToProto(createdBoard),
	}, nil
}
