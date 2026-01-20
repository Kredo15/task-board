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
	// Парсим и валидируем ID доски
	boardID, err := board.ParseBoardID(req.ID)
	if err != nil {
		return nil, err
	}
	// Получаем доску из репозитория
	b, err := uc.repo.GetByID(ctx, boardID)
	if err != nil {
		return nil, err
	}

	var titleToUpdate *board.Title
	if req.Title != nil {
		t, err := board.ParseTitle(*req.Title)
		if err != nil {
			return nil, err
		}
		titleToUpdate = &t
	}

	var descToUpdate *board.Description
	if req.Description != nil {
		d, err := board.ParseDescription(*req.Description)
		if err != nil {
			return nil, err
		}
		descToUpdate = &d
	}

	eventID := board.EventID(uc.gen.Generate())
	// Обновляем доску и генерируем событие
	event, err := b.Update(titleToUpdate, descToUpdate, eventID)
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
