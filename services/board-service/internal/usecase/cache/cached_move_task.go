package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedMoveTaskUseCase struct {
	next  task.MoveTaskUseCase
	cache Invalidator
}

func NewCachedMoveTaskUseCase(next task.MoveTaskUseCase, cache BoardCache) *cachedMoveTaskUseCase {
	return &cachedMoveTaskUseCase{next: next, cache: cache}
}

func (c *cachedMoveTaskUseCase) Execute(ctx context.Context, cmd *task.MoveTaskRequest) (*task.TaskResponse, error) {
	// Сначала выполняем основной сценарий перемещения задачи
	resp, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)
	return resp, nil
}
