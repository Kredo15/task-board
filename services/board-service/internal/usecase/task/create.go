package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/common"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/task"
)

// CreateTaskUseCase представляет обработчик команды создания доски
type CreateTaskUseCase struct {
	repo     task.TaskRepository
	gen      task.IDGenerator
	lexorank common.Ranker
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewCreateTaskUseCase(r task.TaskRepository, g task.IDGenerator) *CreateTaskUseCase {

	return &CreateTaskUseCase{
		repo: r,
		gen:  g,
	}
}

// Execute обрабатывает команду создания задачи
func (uc *CreateTaskUseCase) Execute(ctx context.Context, cmd *CreateTaskRequest) (*TaskResponse, error) {
	colID, err := column.NewColumnID(cmd.ColumnID)
	if err != nil {
		return nil, err
	}
	ranks, err := uc.repo.GetRanksByColumn(ctx, colID)
	if err != nil {
		return nil, err
	}
	// Считаем куда вставлять
	newrank, err := uc.lexorank.CalculateRank(int(cmd.Position), ranks)
	// Преобразование запроса в доменную модель
	newTask, err := task.NewTask(
		uc.gen,
		cmd.ColumnID,
		cmd.Title,
		cmd.Description,
		newrank,
		cmd.AssigneeID,
	)

	if err != nil {
		return nil, err
	}

	// Сохранение доски в репозитории
	if err := uc.repo.Create(ctx, newTask); err != nil {
		return nil, err
	}

	finalIndex, err := uc.repo.GetIndexByRank(ctx, colID, task.Rank(newTask.Rank()))
	if err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &TaskResponse{
		ID:          newTask.ID(),
		ColumnID:    newTask.ColumnID(),
		Title:       newTask.Title(),
		Description: newTask.Description(),
		Position:    finalIndex,
		AssigneeID:  newTask.AssigneeID(),
		CreatedAt:   newTask.CreatedAt(),
		UpdatedAt:   newTask.UpdatedAt(),
	}

	return response, nil
}
