package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

// CreateTaskUseCase представляет обработчик команды создания доски
type DeleteTaskUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewDeleteTaskUseCase(r board.BoardRepository, g board.IDGenerator) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		repo: r,
		gen:  g,
	}
}

// Execute обрабатывает команду удаления задачи
func (uc *DeleteTaskUseCase) Execute(ctx context.Context, req *DeleteTaskRequest) error {
	// Преобразование запроса в доменную модель
	taskID, err := board.ParseTaskID(req.ID)
	if err != nil {
		return err
	}

	boardID, err := board.ParseBoardID(req.BoardID)
	if err != nil {
		return err
	}
	// Загружаем задачу из репозитория
	t, err := uc.repo.GetTaskByID(ctx, boardID, taskID)
	if err != nil {
		return err
	}

	eventID := board.EventID(uc.gen.Generate())
	// Создаем событие удаления задачи
	event := t.Delete(eventID)

	// Удаляем задачу из репозитория и сохраняем событие
	err = uc.repo.DeleteTask(ctx, boardID, taskID, event)
	if err != nil {
		return err
	}

	return nil
}
