package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/task"
)

type MoveTaskUseCase struct {
	repo     task.TaskRepository
	lexorank domain.LexorankGen
}

func NewMoveTaskUseCase(r task.TaskRepository) *MoveTaskUseCase {
	return &MoveTaskUseCase{
		repo: r,
	}
}

func (uc *MoveTaskUseCase) Execute(ctx context.Context, cmd *MoveTaskRequest) (*TaskResponse, error) {
	// Валидируем columnId
	toColumnID, err := column.NewColumnID(cmd.ColumnID)
	if err != nil {
		return nil, err
	}
	// Получаем rank по id tasks
	ranks, err := uc.repo.GetRanksByID(ctx, cmd.AfterTaskID, cmd.BeforeTaskID)
	if err != nil {
		return nil, err
	}
	// Считаем куда вставлять
	newrank, err := uc.lexorank.Between(ranks[0], ranks[1])
	if err != nil {
		return nil, err
	}
	// Валидируем taskId
	task_id, err := task.NewTaskID(cmd.ID)
	if err != nil {
		return nil, err
	}
	// Получаем нужную таску
	t, err := uc.repo.GetByID(ctx, task_id)
	if err != nil {
		return nil, err
	}
	// Валидируем rank
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

	return &TaskResponse{
		t.ID(),
		t.ColumnID(),
		t.Title(),
		t.Description(),
		t.Rank(),
		t.AssigneeID(),
		t.CreatedAt(),
		t.UpdatedAt(),
	}, nil
}
