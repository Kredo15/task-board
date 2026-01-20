package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
)

type cachedAddTaskUseCase struct {
	next  task.AddTaskUseCase
	cache Invalidator
}

func NewCachedAddTaskUseCase(next task.AddTaskUseCase, cache BoardCache) *cachedAddTaskUseCase {
	return &cachedAddTaskUseCase{next: next, cache: cache}
}

func (c *cachedAddTaskUseCase) Execute(ctx context.Context, cmd *task.CreateTaskRequest) (*task.TaskResponse, error) {
	// Сначала выполняем основной сценарий создания задачи
	task, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)
	return task, nil
}
