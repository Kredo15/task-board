package board

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

// createBoardHandler представляет обработчик команды создания доски
type CreateBoardUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewCreateBoardUseCase(r board.BoardRepository, g board.IDGenerator) *CreateBoardUseCase {
	return &CreateBoardUseCase{
		repo: r,
		gen:  g,
	}
}

// Execute обрабатывает команду создания доски
func (h *CreateBoardUseCase) Execute(ctx context.Context, cmd *CreateBoardRequest) (*BoardResponse, error) {
	// Преобразование запроса в доменную модель
	newBoard, event, err := board.NewBoard(
		h.gen,
		cmd.Title,
		cmd.Description,
		cmd.OwnerID,
	)

	if err != nil {
		return nil, err
	}

	// Сохраненяем доску в репозитории и сохраняем событие
	if err := h.repo.Create(ctx, newBoard, event); err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &BoardResponse{
		ID:          newBoard.ID(),
		Title:       newBoard.Title(),
		Description: newBoard.Description(),
		OwnerID:     newBoard.OwnerID(),
		CreatedAt:   newBoard.CreatedAt(),
		UpdatedAt:   newBoard.UpdatedAt(),
	}

	return response, nil
}
