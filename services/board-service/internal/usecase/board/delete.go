package board

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type DeleteBoardUseCase struct {
	repo board.BoardRepository
}

func NewDeleteBoardUseCase(r board.BoardRepository) *DeleteBoardUseCase {
	return &DeleteBoardUseCase{
		repo: r,
	}
}

func (h *DeleteBoardUseCase) Execute(ctx context.Context, cmd *DeleteBoardRequest) error {
	// Преобразование запроса в доменную модель
	boardID, err := board.NewBoardID(cmd.ID)
	if err != nil {
		return err
	}
	// Удаляем доску из репозитория
	err = h.repo.Delete(ctx, boardID)
	if err != nil {
		return err
	}
	return nil
}
