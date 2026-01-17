package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedDeleteTaskUseCase struct {
	next  task.DeleteTaskUseCase
	cache Invalidator
}

func NewCachedDeleteTaskUseCase(next task.DeleteTaskUseCase, cache BoardCache) *cachedDeleteTaskUseCase {
	return &cachedDeleteTaskUseCase{next: next, cache: cache}
}

func (c *cachedDeleteTaskUseCase) Execute(ctx context.Context, cmd *task.DeleteTaskRequest) error {
	// Сначала выполняем основной сценарий удаления задачи
	err := c.next.Execute(ctx, cmd)
	if err != nil {
		return err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)
	return nil
}
