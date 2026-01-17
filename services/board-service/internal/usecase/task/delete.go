package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/task"
)

// CreateTaskUseCase представляет обработчик команды создания доски
type DeleteTaskUseCase struct {
	repo task.TaskRepository
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewDeleteTaskUseCase(r task.TaskRepository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		repo: r,
	}
}

// Execute обрабатывает команду удаления задачи
func (h *DeleteTaskUseCase) Execute(ctx context.Context, cmd *DeleteTaskRequest) error {
	// Преобразование запроса в доменную модель
	taskID, err := task.NewTaskID(cmd.ID)
	if err != nil {
		return err
	}

	err = h.repo.Delete(ctx, taskID)
	if err != nil {
		return err
	}

	return nil
}
