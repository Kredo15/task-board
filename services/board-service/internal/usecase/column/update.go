package column

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

type UpdateColumnUseCase struct {
	repo column.ColumnRepository
}

func NewUpdateTaskUseCase(r column.ColumnRepository) *UpdateColumnUseCase {
	return &UpdateColumnUseCase{
		repo: r,
	}
}

func (uc *UpdateColumnUseCase) Execute(ctx context.Context, cmd *UpdateColumnRequest) (*ColumnResponse, error) {

	// Валидируем column_id
	column_id, err := column.NewColumnID(cmd.ID)
	if err != nil {
		return nil, err
	}
	// Получаем нужную таску
	c, err := uc.repo.GetByID(ctx, column_id)
	if err != nil {
		return nil, err
	}
	// Валидируем title
	title, err := column.NewTitle(cmd.Title)
	if err != nil {
		return nil, err
	}

	// Обновляем позицию и колонку
	c.Update(title)

	// Обновляем сущность
	err_update := uc.repo.Update(ctx, c)
	if err_update != nil {
		return nil, err
	}

	// Сохраняем в транзакции через Outbox
	//err = uc.repo.UpdateAndNotify(ctx, task)

	// Инвалидация кэша Redis
	//uc.cache.Delete(ctx, "board:"+task.BoardID)
	// Так как нам нужно отдавать Position int, смотрим какой по счету текущий rank

	return &ColumnResponse{
		c.ID(),
		c.BoardID(),
		c.Title(),
		c.Rank(),
		c.CreatedAt(),
		c.UpdatedAt(),
	}, nil
}
