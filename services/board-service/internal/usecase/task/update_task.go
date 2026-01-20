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
	taskID, err := board.ParseTaskID(req.ID)
	if err != nil {
		return nil, err
	}

	boardID, err := board.ParseBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}

	// Получаем нужную таску
	t, err := uc.repo.GetTaskByID(ctx, boardID, taskID)
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

	var assigneeIDToUpdate *board.AssigneeID
	if req.AssigneeID != nil {
		assigneeID, err := board.ParseAssigneeID(*req.AssigneeID)
		if err != nil {
			return nil, err
		}
		assigneeIDToUpdate = &assigneeID
	}

	// Генерируем ID события
	eventID := board.EventID(uc.gen.Generate())

	// Обновляем задачу и получаем событие
	event := t.Update(titleToUpdate, descToUpdate, assigneeIDToUpdate, eventID)

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
