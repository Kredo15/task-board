package board

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type UpdateBoardUseCase struct {
	repo board.BoardRepository
}

func NewUpdateBoardUseCase(r board.BoardRepository) *UpdateBoardUseCase {
	return &UpdateBoardUseCase{
		repo: r,
	}
}

func (h *UpdateBoardUseCase) Execute(ctx context.Context, cmd *UpdateBoardRequest) (*BoardResponse, error) {
	// Преобразование запроса в доменную модель
	boardID, err := board.NewBoardID(cmd.ID)
	if err != nil {
		return nil, err
	}
	// Получаем доску из репозитория
	b, err := h.repo.GetByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	// Обновляем доску
	b.Update(cmd.Title, cmd.Description)
	// Сохраняем обновленную доску в репозитории
	err = h.repo.Update(ctx, b)
	if err != nil {
		return nil, err
	}
	// Возвращаем успешный ответ
	response := &BoardResponse{
		ID:          b.ID(),
		Title:       b.Title(),
		Description: b.Description(),
		OwnerID:     b.OwnerID(),
		CreatedAt:   b.CreatedAt(),
	}
	return response, nil
}
