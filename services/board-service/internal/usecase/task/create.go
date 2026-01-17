package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/task"
)

// CreateTaskUseCase представляет обработчик команды создания доски
type CreateTaskUseCase struct {
	repo     task.TaskRepository
	gen      task.IDGenerator
	lexorank domain.LexorankGen
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewCreateTaskUseCase(r task.TaskRepository, g task.IDGenerator, lrank domain.LexorankGen) *CreateTaskUseCase {

	return &CreateTaskUseCase{
		repo:     r,
		gen:      g,
		lexorank: lrank,
	}
}

// Execute обрабатывает команду создания задачи
func (uc *CreateTaskUseCase) Execute(ctx context.Context, cmd *CreateTaskRequest) (*TaskResponse, error) {
	// Валидируем columnId
	colID, err := column.NewColumnID(cmd.ColumnID)
	if err != nil {
		return nil, err
	}
	// Получаем max rank для вставки новой задачи в конец
	rank, err := uc.repo.GetMaxRanksByColumn(ctx, colID)
	if err != nil {
		return nil, err
	}
	// Считаем новый rank
	newrank, err := uc.lexorank.Between(rank, "")
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

	// Возвращаем успешный ответ
	response := &TaskResponse{
		ID:          newTask.ID(),
		ColumnID:    newTask.ColumnID(),
		Title:       newTask.Title(),
		Description: newTask.Description(),
		Rank:        newTask.Rank(),
		AssigneeID:  newTask.AssigneeID(),
		CreatedAt:   newTask.CreatedAt(),
		UpdatedAt:   newTask.UpdatedAt(),
	}

	return response, nil
}
