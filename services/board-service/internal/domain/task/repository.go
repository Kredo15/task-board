package task

import (
	"context"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, id TaskID) error
}
