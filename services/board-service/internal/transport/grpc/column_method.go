package grpc

import (
	"context"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
)

func (h *Handler) AddColumn(ctx context.Context, req *boardv1.CreateColumnRequest) (*boardv1.CreateColumnResponse, error) {
	// Раелизовать метод Add
	return nil, nil
}

func (h *Handler) UpdeteColumn(ctx context.Context, req *boardv1.UpdateColumnRequest) (*boardv1.UpdateColumnResponse, error) {
	// Раелизовать метод Update
	return nil, nil
}

func (h *Handler) MoveColumn(ctx context.Context, req *boardv1.MoveColumnRequest) (*boardv1.MoveColumnResponse, error) {
	// Раелизовать метод Move
	return nil, nil
}

func (h *Handler) DeleteColumn(ctx context.Context, req *boardv1.DeleteColumnRequest) (*boardv1.DeleteColumnResponse, error) {
	// Раелизовать метод Delete
	return nil, nil
}
