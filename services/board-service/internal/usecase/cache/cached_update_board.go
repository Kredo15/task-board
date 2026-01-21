package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

type cachedUpdateBoardUseCase struct {
	next  board.UpdateBoard
	cache Invalidator
}

func NewCachedUpdateBoardUseCase(next board.UpdateBoard, cache BoardCache) *cachedUpdateBoardUseCase {
	return &cachedUpdateBoardUseCase{next: next, cache: cache}
}

func (c *cachedUpdateBoardUseCase) Execute(ctx context.Context, req *board.UpdateBoardRequest) (*board.BoardResponse, error) {
	// Сначала выполняем основной сценарий обновления доски
	resp, err := c.next.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	// сбрасываем кэш
	_ = c.cache.Invalidate(ctx, req.ID)

	return resp, nil
}
