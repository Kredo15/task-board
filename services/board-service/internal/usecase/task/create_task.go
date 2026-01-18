package task

import (
	"context"
	"fmt"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

// CreateTaskUseCase представляет обработчик команды создания доски
type CreateTaskUseCase struct {
	repo     board.BoardRepository
	gen      board.IDGenerator
	lexorank board.LexorankGen
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewCreateTaskUseCase(r board.BoardRepository, g board.IDGenerator, lrank board.LexorankGen) *CreateTaskUseCase {

	return &CreateTaskUseCase{
		repo:     r,
		gen:      g,
		lexorank: lrank,
	}
}

// Execute обрабатывает команду создания задачи
func (uc *CreateTaskUseCase) Execute(ctx context.Context, req *CreateTaskRequest) (*TaskResponse, error) {
	// Валидируем columnId
	colID, err := board.NewColumnID(req.ColumnID)
	if err != nil {
		return nil, err
	}
	boardID, err := board.NewBoardID(req.BoardID)
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

	task, event, err := b.AddTask(
		uc.gen,
		currentCount,
		req.ColumnID,
		req.Title,
		req.Description,
		newrank,
		req.AssigneeID,
	)

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
