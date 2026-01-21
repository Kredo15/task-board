package grpc

import (
	"context"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
)

func (h *Handler) AddTask(ctx context.Context, req *boardv1.CreateTaskRequest) (*boardv1.CreateTaskResponse, error) {
	// Раелизовать метод Add
	return nil, nil
}

func (h *Handler) UpdeteTask(ctx context.Context, req *boardv1.UpdateTaskRequest) (*boardv1.UpdateTaskResponse, error) {
	// Раелизовать метод Update
	return nil, nil
}

func (h *Handler) MoveTask(ctx context.Context, req *boardv1.MoveTaskRequest) (*boardv1.MoveTaskResponse, error) {
	// Раелизовать метод Move
	return nil, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *boardv1.DeleteTaskRequest) (*boardv1.DeleteTaskResponse, error) {
	// Раелизовать метод Delete
	return nil, nil
}
