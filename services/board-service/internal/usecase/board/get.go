package board

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type GetBoardUseCase struct {
	repo board.BoardRepository
}

func NewGetBoardUseCase(r board.BoardRepository) *GetBoardUseCase {
	return &GetBoardUseCase{
		repo: r,
	}
}

func (h *GetBoardUseCase) Execute(ctx context.Context, cmd *GetBoardRequest) (*BoardResponse, error) {
	// Преобразование запроса в доменную модель
	boardID, err := board.NewBoardID(cmd.ID)
	if err != nil {
		return nil, err
	}
	// Получаем доску из репозитория
	board, err := h.repo.GetByID(ctx, boardID)
	if err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &BoardResponse{
		ID:          board.ID(),
		Title:       board.Title(),
		Description: board.Description(),
		OwnerID:     board.OwnerID(),
		CreatedAt:   board.CreatedAt(),
	}

	return response, nil
}
