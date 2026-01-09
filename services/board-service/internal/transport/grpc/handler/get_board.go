package grpc

import (
	"context"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	usecase "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
	"github.com/Kredo15/task-board/services/board-service/pkg/grpcutil"
)

type GetBoard interface {
	Execute(ctx context.Context, cmd *usecase.GetBoardRequest) (*usecase.BoardResponse, error)
}

func (h *BoardHandler) Get(
	ctx context.Context,
	req *boardv1.GetBoardRequest,
) (*boardv1.GetBoardResponse, error) {

	_, err := grpcutil.GetUserID(ctx)
	if err != nil {
		return nil, mapErrorToGRPC(err)
	}

	boardDTO := usecase.GetBoardRequest{
		ID: req.BoardId,
	}

	if err := h.validate.Struct(boardDTO); err != nil {
		h.log.Error("invalid input", err)
		return nil, mapErrorToGRPC(err)
	}

	createdBoard, err := h.getUC.Execute(ctx, &boardDTO)
	if err != nil {
		h.log.Error("createBoard usecase error", err)
		return nil, mapErrorToGRPC(err)
	}

	h.log.Info("board retrieved successfully %s", createdBoard.ID)
	return &boardv1.GetBoardResponse{
		View: mapDTOBoardViewToProto(createdBoard),
	}, nil
}
