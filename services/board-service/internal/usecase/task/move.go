package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/common"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/task"
)

type MoveTaskUseCase struct {
	repo     task.TaskRepository
	lexorank common.Ranker
}

func NewMoveTaskUseCase(r task.TaskRepository) *MoveTaskUseCase {
	return &MoveTaskUseCase{
		repo: r,
	}
}

func (uc *MoveTaskUseCase) Execute(ctx context.Context, cmd *MoveTaskRequest) (*TaskResponse, error) {
	// Получаем все существующие позиции в целевой колонке
	toColumnID, err := column.NewColumnID(cmd.ColumnID)
	if err != nil {
		return nil, err
	}

	ranks, err := uc.repo.GetRanksByColumn(ctx, toColumnID)
	if err != nil {
		return nil, err
	}
	// Считаем куда вставлять
	newrank, err := uc.lexorank.CalculateRank(int(cmd.Position), ranks)
	if err != nil {
		return nil, err
	}

	// Получаем нужную таску
	task_id, err := task.NewTaskID(cmd.ID)
	if err != nil {
		return nil, err
	}

	t, err := uc.repo.GetByID(ctx, task_id)
	if err != nil {
		return nil, err
	}

	toRank, err := task.NewRank(newrank)
	if err != nil {
		return nil, err
	}
	// Обновляем позицию и колонку
	t.Move(toColumnID, toRank)

	// Обновляем сущность
	err_update := uc.repo.Update(ctx, t)
	if err_update != nil {
		return nil, err
	}

	// Сохраняем в транзакции через Outbox
	//err = uc.repo.UpdateAndNotify(ctx, task)

	// Инвалидация кэша Redis
	//uc.cache.Delete(ctx, "board:"+task.BoardID)
	// Так как нам нужно отдавать Position int, смотрим какой по счету текущий rank

	finalIndex, err := uc.repo.GetIndexByRank(ctx, column.ColumnID(t.ColumnID()), toRank)
	if err != nil {
		return nil, err
	}

	return &TaskResponse{
		t.ID(),
		t.ColumnID(),
		t.Title(),
		t.Description(),
		finalIndex,
		t.AssigneeID(),
		t.CreatedAt(),
		t.UpdatedAt(),
	}, nil
}
