package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedCreateColumnUseCase struct {
	next  column.AddColumn
	cache Invalidator
}

func NewCachedCreateColumnUseCase(next column.AddColumn, cache Invalidator) *cachedCreateColumnUseCase {
	return &cachedCreateColumnUseCase{next: next, cache: cache}
}

func (c *cachedCreateColumnUseCase) Execute(ctx context.Context, req *column.AddColumnRequest) (*column.ColumnResponse, error) {
	// Сначала выполняем основной сценарий создания колонки
	resp, err := c.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)

	return resp, nil
}
