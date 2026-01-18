package column

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type UpdateColumnUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

func NewUpdateColumnUseCase(r board.BoardRepository, g board.IDGenerator) *UpdateColumnUseCase {
	return &UpdateColumnUseCase{
		repo: r,
		gen:  g,
	}
}

func (uc *UpdateColumnUseCase) Execute(ctx context.Context, req *UpdateColumnRequest) (*ColumnResponse, error) {

	// Валидируем column_id
	columnID, err := board.NewColumnID(req.ID)
	if err != nil {
		return nil, err
	}
	boardID, err := board.NewBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}
	// Получаем нужную колонку
	c, err := uc.repo.GetColumnByID(ctx, boardID, columnID)
	if err != nil {
		return nil, err
	}

	// Обновляем позицию и колонку
	event, err := c.Update(uc.gen, req.Title)
	if err != nil {
		return nil, err
	}

	// Сохраняем обновленную колонку в репозитории и сохраняем событие
	err = uc.repo.SaveColumn(ctx, boardID, c, event)
	if err != nil {
		return nil, err
	}
	// Возвращаем успешный ответ
	return &ColumnResponse{
		c.ID(),
		c.BoardID(),
		c.Title(),
		c.Rank(),
		c.CreatedAt(),
		c.UpdatedAt(),
	}, nil
}
