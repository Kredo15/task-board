package column

import (
	"context"
	"fmt"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type MoveColumnUseCase struct {
	repo     board.BoardRepository
	gen      board.IDGenerator
	lexorank board.LexorankGen
}

func NewMoveColumnUseCase(r board.BoardRepository, g board.IDGenerator, l board.LexorankGen) *MoveColumnUseCase {
	return &MoveColumnUseCase{
		repo:     r,
		gen:      g,
		lexorank: l,
	}
}

func (uc *MoveColumnUseCase) Execute(ctx context.Context, req *MoveColumnRequest) (*ColumnResponse, error) {
	// Валидируем boardId
	boardID, err := board.NewBoardID(req.BoardID)
	if err != nil {
		return nil, err
	}
	// Валидируем columnId
	columnID, err := board.NewColumnID(req.ID)
	if err != nil {
		return nil, err
	}
	// Получаем rank по id columns
	ranks, err := uc.repo.GetColumnRanks(ctx, boardID)
	if err != nil {
		return nil, err
	}
	// Ищем ранги соседних колонок
	beforeRank, afterRank := uc.findNeighborRanks(req, ranks)
	// Определяем границы для нового ранга
	newrank, err := uc.lexorank.Between(beforeRank, afterRank)
	if err != nil {
		return nil, err
	}
	// Получаем нужную колонку
	c, err := uc.repo.GetColumnByID(ctx, boardID, columnID)
	if err != nil {
		return nil, err
	}
	// Валидируем rank
	toRank, err := board.NewRank(newrank)
	if err != nil {
		return nil, err
	}
	// Обновляем позицию и колонку
	event, err := c.Move(uc.gen, toRank, boardID)
	if err != nil {
		return nil, fmt.Errorf("domain column move: %w", err)
	}
	// Обновляем сущность и сохраняем событие
	err_update := uc.repo.SaveColumn(ctx, boardID, c, event)
	if err_update != nil {
		return nil, err
	}

	return &ColumnResponse{
		c.ID(),
		c.BoardID(),
		c.Title(),
		c.Rank(),
		c.CreatedAt(),
		c.UpdatedAt(),
	}, nil
}

func (uc *MoveColumnUseCase) findNeighborRanks(req *MoveColumnRequest, allColumn []board.ColumnRank) (string, string) {
	var beforeRank, afterRank string
	// Создаем карту для быстрого поиска ранга по ID
	rankMap := make(map[string]string)
	for _, c := range allColumn {
		rankMap[string(c.ID)] = c.Rank
	}

	// Если PrevColumnID указан, берем его ранг из карты
	if req.BeforeColumnID != nil {
		if r, ok := rankMap[*req.BeforeColumnID]; ok {
			beforeRank = r
		}
	}

	// Если NextColumnID указан, берем его ранг из карты
	if req.AfterColumnID != nil {
		if r, ok := rankMap[*req.AfterColumnID]; ok {
			afterRank = r
		}
	}

	return beforeRank, afterRank
}
