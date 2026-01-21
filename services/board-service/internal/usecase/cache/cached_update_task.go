package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedUpdateTaskUseCase struct {
	next  task.UpdateTask
	cache Invalidator
}

func NewCachedUpdateTaskUseCase(next task.UpdateTask, cache BoardCache) *cachedUpdateTaskUseCase {
	return &cachedUpdateTaskUseCase{next: next, cache: cache}
}

func (c *cachedUpdateTaskUseCase) Execute(ctx context.Context, req *task.UpdateTaskRequest) (*task.TaskResponse, error) {
	// Сначала выполняем основной сценарий обновления задачи
	resp, err := c.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)
	return resp, nil
}
