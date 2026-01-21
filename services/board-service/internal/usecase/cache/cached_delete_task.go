package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedDeleteTaskUseCase struct {
	next  task.DeleteTask
	cache Invalidator
}

func NewCachedDeleteTaskUseCase(next task.DeleteTask, cache BoardCache) *cachedDeleteTaskUseCase {
	return &cachedDeleteTaskUseCase{next: next, cache: cache}
}

func (c *cachedDeleteTaskUseCase) Execute(ctx context.Context, req *task.DeleteTaskRequest) error {
	// Сначала выполняем основной сценарий удаления задачи
	err := c.next.Execute(ctx, req)
	if err != nil {
		return err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)
	return nil
}
