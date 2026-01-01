package board

import (
	"context"
)

type BoardRepository interface {
	Create(ctx context.Context, board *Board) error
	GetBoard(ctx context.Context, id BoardID) (*Board, error)
	GetBoards(ctx context.Context) ([]*Board, error)
	Update(ctx context.Context, board Board) error
	Delete(ctx context.Context, id BoardID) error
}
