package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

type cachedGetBoardUseCase struct {
	next  GetBoardUC
	cache BoardCache
}

func NewCachedGetBoardUseCase(next GetBoardUC, cache BoardCache) *cachedGetBoardUseCase {
	return &cachedGetBoardUseCase{next: next, cache: cache}
}

func (c *cachedGetBoardUseCase) Execute(ctx context.Context, cmd *board.GetBoardRequest) (*board.BoardResponse, error) {
	// Пытаемся достать из кэша по ID из команды
	dto, err := c.cache.Get(ctx, cmd.ID)
	if err == nil {
		return dto, nil
	}

	// Если нет в кэше — выполняем реальный сценарий
	resp, err := c.next.Execute(ctx, cmd)

	// Сохраняем результат в кэш
	if err == nil {
		c.cache.Set(ctx, resp)
	}

	return resp, err
}
