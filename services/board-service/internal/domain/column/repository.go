package column

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type ColumnRepository interface {
	Create(ctx context.Context, column *Column) error
	GetByID(ctx context.Context, id ColumnID) (*Column, error)
	Update(ctx context.Context, column *Column) error
	GetByBoard(ctx context.Context, boardID board.BoardID) ([]*Column, error)
	GetByBoards(ctx context.Context, boardsID []board.BoardID) ([]*Column, error)
	// получение rank по id column
	GetRanksByID(ctx context.Context, afterColumnID, beforeColumnID string) ([]string, error)
	// получение max rank для добавления новой колонки в конец
	GetMaxRanksByBoard(ctx context.Context, boardID board.BoardID) (string, error)
	Delete(ctx context.Context, id ColumnID) error
}
