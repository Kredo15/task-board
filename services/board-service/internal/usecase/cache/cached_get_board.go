package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

type cachedGetBoardUseCase struct {
	next  board.GetBoard
	cache BoardCache
}

func NewCachedGetBoardUseCase(next board.GetBoard, cache BoardCache) *cachedGetBoardUseCase {
	return &cachedGetBoardUseCase{next: next, cache: cache}
}

func (c *cachedGetBoardUseCase) Execute(ctx context.Context, req *board.GetBoardRequest) (*board.BoardResponse, error) {
	// Пытаемся достать из кэша по ID из команды
	dto, err := c.cache.Get(ctx, req.ID)
	if err == nil {
		return dto, nil
	}

	// Если нет в кэше — выполняем реальный сценарий
	resp, err := c.next.Execute(ctx, req)

	// Сохраняем результат в кэш
	if err == nil {
		c.cache.Set(ctx, resp)
	}

	return resp, err
}
