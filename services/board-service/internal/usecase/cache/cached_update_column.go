package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedUpdateColumnUseCase struct {
	next  column.UpdateColumnUseCase
	cache Invalidator
}

func NewCachedUpdateColumnUseCase(next column.UpdateColumnUseCase, cache BoardCache) *cachedUpdateColumnUseCase {
	return &cachedUpdateColumnUseCase{next: next, cache: cache}
}

func (c *cachedUpdateColumnUseCase) Execute(ctx context.Context, cmd *column.UpdateColumnRequest) (*column.ColumnResponse, error) {
	// Сначала выполняем основной сценарий обновления колонки
	column, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)
	return column, nil
}
