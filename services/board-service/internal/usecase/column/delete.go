package column

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

// DeleteColumnUseCase представляет обработчик команды удаления колонки
type DeleteColumnUseCase struct {
	repo column.ColumnRepository
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewDeleteTaskUseCase(r column.ColumnRepository) *DeleteColumnUseCase {
	return &DeleteColumnUseCase{
		repo: r,
	}
}

// Execute обрабатывает команду удаления задачи
func (h *DeleteColumnUseCase) Execute(ctx context.Context, cmd *DeleteColumnRequest) (*DeleteColumnResponse, error) {
	// Преобразование запроса в доменную модель
	columnID, err := column.NewColumnID(cmd.ID)
	if err != nil {
		return nil, err
	}

	err = h.repo.Delete(ctx, columnID)
	if err != nil {
		return nil, err
	}

	return &DeleteColumnResponse{}, nil
}
