package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedCreateColumnUseCase struct {
	next  column.CreateColumnUseCase
	cache Invalidator
}

func NewCachedCreateColumnUseCase(next column.CreateColumnUseCase, cache BoardCache) *cachedCreateColumnUseCase {
	return &cachedCreateColumnUseCase{next: next, cache: cache}
}

func (c *cachedCreateColumnUseCase) Execute(ctx context.Context, cmd *column.CreateColumnRequest) (*column.ColumnResponse, error) {
	// Сначала выполняем основной сценарий создания колонки
	resp, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)

	return resp, nil
}
