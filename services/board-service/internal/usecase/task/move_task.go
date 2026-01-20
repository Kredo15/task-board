package task

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type MoveTaskUseCase struct {
	repo     board.BoardRepository
	gen      board.IDGenerator
	lexorank board.LexorankGen
}

func NewMoveTaskUseCase(r board.BoardRepository, g board.IDGenerator, l board.LexorankGen) *MoveTaskUseCase {
	return &MoveTaskUseCase{
		repo:     r,
		gen:      g,
		lexorank: l,
	}
}

func (uc *MoveTaskUseCase) Execute(ctx context.Context, req *MoveTaskRequest) (*TaskResponse, error) {
	// Валидируем columnId
	toColumnID, err := board.ParseColumnID(req.ColumnID)
	if err != nil {
		return nil, err
	}
	boardID, err := board.ParseBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}

	// Получаем rank по id tasks
	ranks, err := uc.repo.GetTaskRanks(ctx, boardID, toColumnID)
	if err != nil {
		return nil, err
	}
	// Ищем ранги соседних тасок
	beforeRank, afterRank := uc.findNeighborRanks(req, ranks)

	// Определяем границы для нового ранга
	newrank, err := uc.lexorank.Between(beforeRank, afterRank)
	if err != nil {
		return nil, err
	}
	// Валидируем taskId
	task_id, err := board.ParseTaskID(req.ID)
	if err != nil {
		return nil, err
	}
	// Получаем нужную таску
	t, err := uc.repo.GetTaskByID(ctx, boardID, task_id)
	if err != nil {
		return nil, err
	}
	// Валидируем rank
	toRank, err := board.ParseRank(newrank)
	if err != nil {
		return nil, err
	}

	// Генерируем ID события
	eventID := board.EventID(uc.gen.Generate())

	// Обновляем позицию и колонку, получаем событие
	event := t.Move(toColumnID, toRank, eventID)

	// Обновляем сущность Task в репозитории и сохраняем событие
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

func (uc *MoveTaskUseCase) findNeighborRanks(req *MoveTaskRequest, allTasks []board.TaskRank) (string, string) {
	var beforeRank, afterRank string
	// Создаем карту для быстрого поиска ранга по ID
	rankMap := make(map[string]string)
	for _, t := range allTasks {
		rankMap[string(t.ID)] = t.Rank
	}

	// Если PrevTaskID указан, берем его ранг из карты
	if req.BeforeTaskID != nil {
		if r, ok := rankMap[*req.BeforeTaskID]; ok {
			beforeRank = r
		}
	}

	// Если NextTaskID указан, берем его ранг из карты
	if req.AfterTaskID != nil {
		if r, ok := rankMap[*req.AfterTaskID]; ok {
			afterRank = r
		}
	}

	return beforeRank, afterRank
}
