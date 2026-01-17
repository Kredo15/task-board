package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
)

type cachedDeleteColumnUseCase struct {
	next  column.DeleteColumnUseCase
	cache Invalidator
}

func NewCachedDeleteColumnUseCase(next column.DeleteColumnUseCase, cache BoardCache) *cachedDeleteColumnUseCase {
	return &cachedDeleteColumnUseCase{next: next, cache: cache}
}

func (c *cachedDeleteColumnUseCase) Execute(ctx context.Context, cmd *column.DeleteColumnRequest) error {
	// Сначала выполняем основной сценарий удаления колонки
	err := c.next.Execute(ctx, cmd)
	if err != nil {
		return err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.BoardID)

	return nil
}
