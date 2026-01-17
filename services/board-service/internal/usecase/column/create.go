package column

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

// CreateColumnUseCase представляет обработчик команды создания колонки
type CreateColumnUseCase struct {
	repo     column.ColumnRepository
	gen      column.IDGenerator
	lexorank domain.LexorankGen
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewCreateColumnUseCase(r column.ColumnRepository, g column.IDGenerator, lrank domain.LexorankGen) *CreateColumnUseCase {

	return &CreateColumnUseCase{
		repo:     r,
		gen:      g,
		lexorank: lrank,
	}
}

// Execute обрабатывает команду создания колонки
func (uc *CreateColumnUseCase) Execute(ctx context.Context, cmd *CreateColumnRequest) (*ColumnResponse, error) {
	// Валидируем columnId
	bID, err := board.NewBoardID(cmd.BoardID)
	if err != nil {
		return nil, err
	}
	// Получаем max rank для вставки новой колонки в конец
	rank, err := uc.repo.GetMaxRanksByBoard(ctx, bID)
	if err != nil {
		return nil, err
	}
	// Считаем новый rank
	newrank, err := uc.lexorank.Between(rank, "")
	// Преобразование запроса в доменную модель
	newTask, err := column.NewColumn(
		uc.gen,
		cmd.BoardID,
		cmd.Title,
		newrank,
	)

	if err != nil {
		return nil, err
	}

	// Сохранение доски в репозитории
	if err := uc.repo.Create(ctx, newTask); err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &ColumnResponse{
		ID:        newTask.ID(),
		BoardID:   newTask.BoardID(),
		Title:     newTask.Title(),
		Rank:      newTask.Rank(),
		CreatedAt: newTask.CreatedAt(),
		UpdatedAt: newTask.UpdatedAt(),
	}

	return response, nil
}
