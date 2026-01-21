package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedUpdateColumnUseCase struct {
	next  column.UpdateColumn
	cache Invalidator
}

func NewCachedUpdateColumnUseCase(next column.UpdateColumn, cache BoardCache) *cachedUpdateColumnUseCase {
	return &cachedUpdateColumnUseCase{next: next, cache: cache}
}

func (c *cachedUpdateColumnUseCase) Execute(ctx context.Context, req *column.UpdateColumnRequest) (*column.ColumnResponse, error) {
	// Сначала выполняем основной сценарий обновления колонки
	column, err := c.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)
	return column, nil
}
