package grpc

import (
	"context"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	usecase "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
	"github.com/Kredo15/task-board/services/board-service/pkg/grpcutil"
)

func (h *Handler) Create(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.CreateBoardResponse, error) {

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

	createdBoard, err := h.board.Create.Execute(ctx, &boardDTO)
	if err != nil {
		h.log.Error("createBoard usecase error", err)
		return nil, mapErrorToGRPC(err)
	}

	h.log.Info("board retrieved successfully %s", createdBoard.ID)
	return &boardv1.CreateBoardResponse{
		Board: mapDTOBoardToProto(createdBoard),
	}, nil
}

func (h *Handler) GetBoard(ctx context.Context, req *boardv1.GetBoardRequest) (*boardv1.GetBoardResponse, error) {

	boardDTO := usecase.GetBoardRequest{
		ID: req.BoardId,
	}

	getBoard, err := h.board.Get.Execute(ctx, &boardDTO)
	if err != nil {
		h.log.Error("createBoard usecase error", err)
		return nil, mapErrorToGRPC(err)
	}

	h.log.Info("board retrieved successfully %s", getBoard.ID)
	return &boardv1.GetBoardResponse{
		View: mapDTOBoardViewToProto(getBoard),
	}, nil
}

func (h *Handler) UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.UpdateBoardResponse, error) {

	boardDTO := usecase.UpdateBoardRequest{
		ID:          req.BoardId,
		Title:       &req.Title,
		Description: &req.Description,
	}

	updateBoard, err := h.board.Update.Execute(ctx, &boardDTO)

	if err != nil {
		h.log.Error("update board usecase error", err)
		return nil, mapErrorToGRPC(err)
	}

	h.log.Info("board updated successfully %s", updateBoard.ID)
	return &boardv1.UpdateBoardResponse{
		Board: mapDTOBoardToProto(updateBoard),
	}, nil

}

func (h *Handler) DeleteBoard(ctx context.Context, req *boardv1.DeleteBoardRequest) (*boardv1.DeleteBoardResponse, error) {

	boardDTO := usecase.DeleteBoardRequest{
		ID: req.BoardId,
	}

	err := h.board.Delete.Execute(ctx, &boardDTO)

	if err != nil {
		h.log.Error("delete board usecase error", err)
		return nil, mapErrorToGRPC(err)
	}

	h.log.Info("delete board successfully %s", req.BoardId)
	return &boardv1.DeleteBoardResponse{}, nil
}
