package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedMoveColumnUseCase struct {
	next  column.MoveColumn
	cache Invalidator
}

func NewCachedMoveColumnUseCase(next column.MoveColumn, cache BoardCache) *cachedMoveColumnUseCase {
	return &cachedMoveColumnUseCase{next: next, cache: cache}
}

func (c *cachedMoveColumnUseCase) Execute(ctx context.Context, req *column.MoveColumnRequest) (*column.ColumnResponse, error) {
	// Сначала выполняем основной сценарий перемещения колонки
	column, err := c.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)
	return column, nil
}
