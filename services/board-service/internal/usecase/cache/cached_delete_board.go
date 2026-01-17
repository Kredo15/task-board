package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

type cachedDeleteBoardUseCase struct {
	next  DeleteBoardUC
	cache Invalidator
}

func NewCachedDeleteBoardUseCase(next DeleteBoardUC, cache Invalidator) *cachedDeleteBoardUseCase {
	return &cachedDeleteBoardUseCase{next: next, cache: cache}
}

func (c *cachedDeleteBoardUseCase) Execute(ctx context.Context, cmd *board.DeleteBoardRequest) error {
	// Сначала удаляем из БД через основной сценарий
	if err := c.next.Execute(ctx, cmd); err != nil {
		return err
	}

	// Если удаление из БД прошло успешно — сбрасываем кэш
	// Возвращать ошибку кэша не нужно, так как данные в БД уже удалены.
	_ = c.cache.Invalidate(ctx, cmd.ID)

	return nil
}
