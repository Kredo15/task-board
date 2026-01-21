package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedDeleteColumnUseCase struct {
	next  column.DeleteColumn
	cache Invalidator
}

func NewCachedDeleteColumnUseCase(next column.DeleteColumn, cache BoardCache) *cachedDeleteColumnUseCase {
	return &cachedDeleteColumnUseCase{next: next, cache: cache}
}

func (c *cachedDeleteColumnUseCase) Execute(ctx context.Context, req *column.DeleteColumnRequest) error {
	// Сначала выполняем основной сценарий удаления колонки
	err := c.next.Execute(ctx, req)
	if err != nil {
		return err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.BoardID)

	return nil
}
