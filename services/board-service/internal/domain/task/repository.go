package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id TaskID) (*Task, error)
	Update(ctx context.Context, task *Task) error
	GetRanksByColumn(ctx context.Context, colID column.ColumnID) ([]string, error)
	GetIndexByRank(ctx context.Context, colID column.ColumnID, rank Rank) (int32, error)
	Delete(ctx context.Context, id TaskID) error
}
