package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedMoveTaskUseCase struct {
	next  task.MoveTask
	cache Invalidator
}

func NewCachedMoveTaskUseCase(next task.MoveTask, cache BoardCache) *cachedMoveTaskUseCase {
	return &cachedMoveTaskUseCase{next: next, cache: cache}
}

func (c *cachedMoveTaskUseCase) Execute(ctx context.Context, req *task.MoveTaskRequest) (*task.TaskResponse, error) {
	// Сначала выполняем основной сценарий перемещения задачи
	resp, err := c.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)
	return resp, nil
}
