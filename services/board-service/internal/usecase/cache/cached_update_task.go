package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedUpdateTaskUseCase struct {
	next  task.UpdateTaskUseCase
	cache Invalidator
}

func NewCachedUpdateTaskUseCase(next task.UpdateTaskUseCase, cache BoardCache) *cachedUpdateTaskUseCase {
	return &cachedUpdateTaskUseCase{next: next, cache: cache}
}

func (c *cachedUpdateTaskUseCase) Execute(ctx context.Context, cmd *task.UpdateTaskRequest) (*task.TaskResponse, error) {
	// Сначала выполняем основной сценарий обновления задачи
	resp, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)
	return resp, nil
}
