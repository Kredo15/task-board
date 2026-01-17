package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/task"
)

type UpdateTaskUseCase struct {
	repo task.TaskRepository
}

func NewUpdateTaskUseCase(r task.TaskRepository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		repo: r,
	}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, cmd *UpdateTaskRequest) (*TaskResponse, error) {

	// Валидируем task_id
	task_id, err := task.NewTaskID(cmd.ID)
	if err != nil {
		return nil, err
	}
	// Получаем нужную таску
	t, err := uc.repo.GetByID(ctx, task_id)
	if err != nil {
		return nil, err
	}
	// Валидируем title
	title, err := task.NewTitle(cmd.Title)
	if err != nil {
		return nil, err
	}
	// Валидируем description
	desc, err := task.NewDescription(cmd.Description)
	if err != nil {
		return nil, err
	}
	// Валидируем description
	aid, err := task.NewAssigneeID(cmd.AssigneeID)
	if err != nil {
		return nil, err
	}
	// Обновляем позицию и колонку
	t.Update(title, desc, aid)

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
