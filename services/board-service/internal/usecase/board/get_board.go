package board

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type GetBoardUseCase struct {
	repo board.BoardRepository
}

func NewGetBoardUseCase(r board.BoardRepository) *GetBoardUseCase {
	return &GetBoardUseCase{
		repo: r,
	}
}

func (h *GetBoardUseCase) Execute(ctx context.Context, cmd *GetBoardRequest) (*BoardResponse, error) {
	// Преобразование запроса в доменную модель
	boardID, err := board.ParseBoardID(cmd.ID)
	if err != nil {
		return nil, err
	}
	// Получаем доску из репозитория
	board, err := h.repo.GetFullBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &BoardResponse{
		ID:          board.ID(),
		Title:       board.Title(),
		Description: board.Description(),
		OwnerID:     board.OwnerID(),
		CreatedAt:   board.CreatedAt(),
		UpdatedAt:   board.UpdatedAt(),
		Columns:     make([]ColumnResponse, 0),
	}

	for _, col := range board.Columns() {
		columnResp := ColumnResponse{
			ID:        col.ID(),
			Title:     col.Title(),
			Rank:      col.Rank(),
			CreatedAt: col.CreatedAt(),
			UpdatedAt: col.UpdatedAt(),
			Tasks:     make([]TaskResponse, 0),
		}
		for _, task := range col.Tasks() {
			taskResp := TaskResponse{
				ID:          task.ID(),
				ColumnID:    task.ColumnID(),
				Title:       task.Title(),
				Description: task.Description(),
				Rank:        task.Rank(),
				AssigneeID:  task.AssigneeID(),
				CreatedAt:   task.CreatedAt(),
				UpdatedAt:   task.UpdatedAt(),
			}
			columnResp.Tasks = append(columnResp.Tasks, taskResp)
		}

		response.Columns = append(response.Columns, columnResp)
	}

	return response, nil
}
