package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedAddTaskUseCase struct {
	next  task.AddTask
	cache Invalidator
}

func NewCachedAddTaskUseCase(next task.AddTask, cache BoardCache) *cachedAddTaskUseCase {
	return &cachedAddTaskUseCase{next: next, cache: cache}
}

func (c *cachedAddTaskUseCase) Execute(ctx context.Context, req *task.CreateTaskRequest) (*task.TaskResponse, error) {
	// Сначала выполняем основной сценарий создания задачи
	task, err := c.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)
	return task, nil
}
