package column

import (
	"context"
	"fmt"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

// DeleteColumnUseCase представляет обработчик команды удаления колонки
type DeleteColumnUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewDeleteTaskUseCase(r board.BoardRepository, gen board.IDGenerator) *DeleteColumnUseCase {
	return &DeleteColumnUseCase{
		repo: r,
		gen:  gen,
	}
}

// Execute обрабатывает команду удаления задачи
func (uc *DeleteColumnUseCase) Execute(ctx context.Context, req *DeleteColumnRequest) error {
	// Преобразование запроса в доменную модель
	columnID, err := board.ParseColumnID(req.ID)
	if err != nil {
		return err
	}
	boardID, err := board.ParseBoardID(req.BoardID)
	if err != nil {
		return err
	}
	// Загружаем колонку из репозитория
	c, err := uc.repo.GetColumnByID(ctx, boardID, columnID)
	if err != nil {
		return fmt.Errorf("failed to get column: %w", err)
	}

	eventID := board.EventID(uc.gen.Generate())
	// Создаем событие удаления колонки
	event := c.Delete(eventID)

	err = uc.repo.DeleteColumn(ctx, boardID, columnID, event)
	if err != nil {
		return err
	}

	return nil
}
