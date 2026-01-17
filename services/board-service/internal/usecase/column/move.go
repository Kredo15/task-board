package column

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

type MoveColumnUseCase struct {
	repo     column.ColumnRepository
	lexorank domain.LexorankGen
}

func NewMoveColumnUseCase(r column.ColumnRepository) *MoveColumnUseCase {
	return &MoveColumnUseCase{
		repo: r,
	}
}

func (uc *MoveColumnUseCase) Execute(ctx context.Context, cmd *MoveColumnRequest) (*ColumnResponse, error) {

	// Получаем rank по id columns
	ranks, err := uc.repo.GetRanksByID(ctx, cmd.AfterColumnID, cmd.BeforeColumnID)
	if err != nil {
		return nil, err
	}
	// Считаем куда вставлять
	newrank, err := uc.lexorank.Between(ranks[0], ranks[1])
	if err != nil {
		return nil, err
	}
	// Валидируем taskId
	column_id, err := column.NewColumnID(cmd.ID)
	if err != nil {
		return nil, err
	}
	// Получаем нужную таску
	c, err := uc.repo.GetByID(ctx, column_id)
	if err != nil {
		return nil, err
	}
	// Валидируем rank
	toRank, err := column.NewRank(newrank)
	if err != nil {
		return nil, err
	}
	// Обновляем позицию и колонку
	c.Move(toRank)
	// Обновляем сущность
	err_update := uc.repo.Update(ctx, c)
	if err_update != nil {
		return nil, err
	}

	// Сохраняем в транзакции через Outbox
	//err = uc.repo.UpdateAndNotify(ctx, task)

	return &ColumnResponse{
		c.ID(),
		c.BoardID(),
		c.Title(),
		c.Rank(),
		c.CreatedAt(),
		c.UpdatedAt(),
	}, nil
}
