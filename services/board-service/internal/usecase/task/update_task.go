package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type UpdateTaskUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

func NewUpdateTaskUseCase(r board.BoardRepository, g board.IDGenerator) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		repo: r,
		gen:  g,
	}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, req *UpdateTaskRequest) (*TaskResponse, error) {

	// Валидируем task_id
	taskID, err := board.NewTaskID(req.ID)
	if err != nil {
		return nil, err
	}

	boardID, err := board.NewBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}

	// Получаем нужную таску
	t, err := uc.repo.GetTaskByID(ctx, boardID, taskID)
	if err != nil {
		return nil, err
	}
	// Обновляем задачу и получаем событие
	event, err := t.Update(uc.gen, req.Title, req.Description, req.AssigneeID)

	// Обновляем сущность в репозитории и сохраняем событие
	err_update := uc.repo.SaveTask(ctx, boardID, t, event)
	if err_update != nil {
		return nil, err
	}

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
