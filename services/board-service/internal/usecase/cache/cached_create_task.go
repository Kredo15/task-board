package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedCreateTaskUseCase struct {
	next  task.CreateTaskUseCase
	cache Invalidator
}

func NewCachedCreateTaskUseCase(next task.CreateTaskUseCase, cache BoardCache) *cachedCreateTaskUseCase {
	return &cachedCreateTaskUseCase{next: next, cache: cache}
}

func (c *cachedCreateTaskUseCase) Execute(ctx context.Context, cmd *task.CreateTaskRequest) (*task.TaskResponse, error) {
	// Сначала выполняем основной сценарий создания задачи
	task, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)
	return task, nil
}
