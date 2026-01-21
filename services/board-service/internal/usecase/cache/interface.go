package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

// Определяем интерфейс для конкретной операции
type BoardCache interface {
	Get(ctx context.Context, id string) (*board.BoardResponse, error)
	Set(ctx context.Context, dto *board.BoardResponse) error
	Invalidator
}

type Invalidator interface {
	Invalidate(ctx context.Context, id string) error
}
