package cache

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

// Определяем интерфейс для конкретной операции
type BoardCache interface {
	Get(ctx context.Context, id string) (*board.BoardResponse, error)
	Set(ctx context.Context, dto *board.BoardResponse)
	Invalidator
}

type Invalidator interface {
	Invalidate(ctx context.Context, id string) error
}

type GetBoardUC interface {
	Execute(ctx context.Context, cmd *board.GetBoardRequest) (*board.BoardResponse, error)
}

type DeleteBoardUC interface {
	Execute(ctx context.Context, cmd *board.DeleteBoardRequest) error
}

type UpdateBoardUC interface {
	Execute(ctx context.Context, cmd *board.UpdateBoardRequest) (*board.BoardResponse, error)
}
