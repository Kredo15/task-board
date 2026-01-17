package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

type cachedUpdateBoardUseCase struct {
	next  UpdateBoardUC
	cache Invalidator
}

func NewCachedUpdateBoardUseCase(next UpdateBoardUC, cache BoardCache) *cachedUpdateBoardUseCase {
	return &cachedUpdateBoardUseCase{next: next, cache: cache}
}

func (c *cachedUpdateBoardUseCase) Execute(ctx context.Context, cmd *board.UpdateBoardRequest) (*board.BoardResponse, error) {
	// Сначала выполняем основной сценарий обновления доски
	resp, err := c.next.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, cmd.ID)

	return resp, nil
}
