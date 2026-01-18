package board

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type UpdateBoardUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

func NewUpdateBoardUseCase(r board.BoardRepository, gen board.IDGenerator) *UpdateBoardUseCase {
	return &UpdateBoardUseCase{
		repo: r,
		gen:  gen,
	}
}

func (uc *UpdateBoardUseCase) Execute(ctx context.Context, req *UpdateBoardRequest) (*BoardResponse, error) {
	// Преобразование запроса в доменную модель
	boardID, err := board.NewBoardID(req.ID)
	if err != nil {
		return nil, err
	}
	// Получаем доску из репозитория
	b, err := uc.repo.GetByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	// Обновляем доску и генерируем событие
	event, err := b.Update(uc.gen, req.Title, req.Description)
	if err != nil {
		return nil, err
	}
	// Сохраняем обновленную доску в репозитории и сохраняем событие
	if err := uc.repo.Update(ctx, b, event); err != nil {
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
