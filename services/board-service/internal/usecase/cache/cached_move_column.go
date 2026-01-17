package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedMoveColumnUseCase struct {
	next  column.MoveColumnUseCase
	cache Invalidator
}

func NewCachedMoveColumnUseCase(next column.MoveColumnUseCase, cache BoardCache) *cachedMoveColumnUseCase {
	return &cachedMoveColumnUseCase{next: next, cache: cache}
}

func (c *cachedMoveColumnUseCase) Execute(ctx context.Context, cmd *column.MoveColumnRequest) (*column.ColumnResponse, error) {
	// Сначала выполняем основной сценарий перемещения колонки
	column, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)
	return column, nil
}
