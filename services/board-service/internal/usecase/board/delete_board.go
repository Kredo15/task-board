package board

import (
	"context"
	"fmt"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type DeleteBoardUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

func NewDeleteBoardUseCase(r board.BoardRepository, gen board.IDGenerator) *DeleteBoardUseCase {
	return &DeleteBoardUseCase{
		repo: r,
		gen:  gen,
	}
}

func (uc *DeleteBoardUseCase) Execute(ctx context.Context, cmd *DeleteBoardRequest) error {
	// Преобразование запроса в доменную модель
	boardID, err := board.NewBoardID(cmd.ID)
	if err != nil {
		return err
	}
	// Загружаем корень агрегата (Доску)
	b, err := uc.repo.GetByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("failed to get board: %w", err)
	}
	// Создаем событие удаления доски
	event, err := b.Delete(uc.gen)
	if err != nil {
		return fmt.Errorf("failed to event delete board: %w", err)
	}
	// Удаляем доску из репозитория
	err = uc.repo.Delete(ctx, boardID, event)
	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}
	return nil
}
