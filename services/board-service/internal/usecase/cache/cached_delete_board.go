package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

type cachedDeleteBoardUseCase struct {
	next  board.DeleteBoard
	cache Invalidator
}

func NewCachedDeleteBoardUseCase(next board.DeleteBoard, cache Invalidator) *cachedDeleteBoardUseCase {
	return &cachedDeleteBoardUseCase{next: next, cache: cache}
}

func (c *cachedDeleteBoardUseCase) Execute(ctx context.Context, req *board.DeleteBoardRequest) error {
	// Сначала удаляем из БД через основной сценарий
	if err := c.next.Execute(ctx, req); err != nil {
		return err
	}

	// Если удаление из БД прошло успешно — сбрасываем кэш
	// Возвращать ошибку кэша не нужно, так как данные в БД уже удалены.
	_ = c.cache.Invalidate(ctx, req.ID)

	return nil
}
