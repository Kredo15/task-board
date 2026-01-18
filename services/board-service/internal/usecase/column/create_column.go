package column

import (
	"context"
	"fmt"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

// AddColumnUseCase представляет обработчик команды создания колонки
type AddColumnUseCase struct {
	repo     board.BoardRepository
	gen      board.IDGenerator
	lexorank board.LexorankGen
}

// NewAddColumnUseCase создает новый экземпляр обработчика команды создания колонки
func NewAddColumnUseCase(r board.BoardRepository, g board.IDGenerator, lrank board.LexorankGen) *AddColumnUseCase {
	return &AddColumnUseCase{
		repo:     r,
		gen:      g,
		lexorank: lrank,
	}
}

// Execute обрабатывает команду создания колонки
func (uc *AddColumnUseCase) Execute(ctx context.Context, req *AddColumnRequest) (*ColumnResponse, error) {
	// Валидируем columnId
	bID, err := board.NewBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}
	// Получаем все ranks для вставки новой колонки в конец
	ranks, err := uc.repo.GetColumnRanks(ctx, bID)
	if err != nil {
		return nil, fmt.Errorf("get column ranks: %w", err)
	}

	var rank string
	if len(ranks) > 0 {
		rank = ranks[len(ranks)-1].Rank
	}
	// Считаем новый rank
	newrank, err := uc.lexorank.Between(rank, "")

	if err != nil {
		return nil, fmt.Errorf("calculate lexorank: %w", err)
	}
	boardID, err := board.NewBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}
	// Загружаем корень агрегата (Доску)
	b, err := uc.repo.GetByID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	column, event, err := b.AddColumn(uc.gen, req.Title, newrank)
	if err != nil {
		return nil, err
	}

	// Сохранение доски в репозитории
	if err := uc.repo.SaveColumn(ctx, bID, column, event); err != nil {
		return nil, fmt.Errorf("save column: %w", err)
	}

	// Возвращаем успешный ответ
	response := &ColumnResponse{
		ID:        column.ID(),
		BoardID:   column.BoardID(),
		Title:     column.Title(),
		Rank:      column.Rank(),
		CreatedAt: column.CreatedAt(),
		UpdatedAt: column.UpdatedAt(),
	}

	return response, nil
}
