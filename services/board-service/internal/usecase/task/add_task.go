package task

import (
	"context"
	"fmt"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

// AddTaskUseCase представляет обработчик команды создания задачи
type AddTaskUseCase struct {
	repo     board.BoardRepository
	gen      board.IDGenerator
	lexorank board.LexorankGen
}

// NewAddTaskUseCase создает новый экземпляр обработчика команды создания задачи
func NewAddTaskUseCase(r board.BoardRepository, g board.IDGenerator, lrank board.LexorankGen) *AddTaskUseCase {

	return &AddTaskUseCase{
		repo:     r,
		gen:      g,
		lexorank: lrank,
	}
}

// Execute обрабатывает команду создания задачи
func (uc *AddTaskUseCase) Execute(ctx context.Context, req *CreateTaskRequest) (*TaskResponse, error) {
	// Валидируем columnId
	colID, err := board.ParseColumnID(req.ColumnID)
	if err != nil {
		return nil, err
	}
	boardID, err := board.ParseBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}
	// Получаем max rank для вставки новой задачи в конец
	ranks, err := uc.repo.GetTaskRanks(ctx, boardID, colID)
	if err != nil {
		return nil, err
	}
	var rank string
	currentCount := len(ranks)
	if currentCount > 0 {
		rank = ranks[currentCount-1].Rank
	}
	// Считаем новый rank
	newrank, err := uc.lexorank.Between(rank, "")
	if err != nil {
		return nil, err
	}

	// Загружаем корень агрегата (Доску)
	b, err := uc.repo.GetByID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	toRank, err := board.ParseRank(newrank)
	if err != nil {
		return nil, err
	}

	title, err := board.ParseTitle(req.Title)
	if err != nil {
		return nil, err
	}

	description, err := board.ParseDescription(req.Description)
	if err != nil {
		return nil, err
	}

	assigneeID, err := board.ParseAssigneeID(req.AssigneeID)
	if err != nil {
		return nil, err
	}

	// Генерируем ID для новой задачи и события
	taskID := board.TaskID(uc.gen.Generate())
	eventID := board.EventID(uc.gen.Generate())

	task, event, err := b.AddTask(taskID, currentCount, colID, title, description, toRank, assigneeID, eventID)

	if err != nil {
		return nil, err
	}

	// Сохранение доски в репозитории
	if err := uc.repo.SaveTask(ctx, boardID, task, event); err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &TaskResponse{
		ID:          task.ID(),
		ColumnID:    task.ColumnID(),
		Title:       task.Title(),
		Description: task.Description(),
		Rank:        task.Rank(),
		AssigneeID:  task.AssigneeID(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
	}

	return response, nil
}
