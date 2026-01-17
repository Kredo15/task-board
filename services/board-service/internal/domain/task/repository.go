package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id TaskID) (*Task, error)
	GetByColumn(ctx context.Context, colID column.ColumnID) ([]*Task, error)
	GetByColumns(ctx context.Context, colsID []column.ColumnID) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id TaskID) error
	// получение rank по id task
	GetRanksByID(ctx context.Context, afterTaskID, beforeTaskID string) ([]string, error)
	// получение max rank для добавления новой задачи в конец
	GetMaxRanksByColumn(ctx context.Context, colID column.ColumnID) (string, error)
}
