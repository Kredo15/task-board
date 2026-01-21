package task

import (
	"context"
)

type AddTask interface {
	Execute(ctx context.Context, req *CreateTaskRequest) (*TaskResponse, error)
}

type MoveTask interface {
	Execute(ctx context.Context, req *MoveTaskRequest) (*TaskResponse, error)
}

type UpdateTask interface {
	Execute(ctx context.Context, req *UpdateTaskRequest) (*TaskResponse, error)
}

type DeleteTask interface {
	Execute(ctx context.Context, req *DeleteTaskRequest) error
}
