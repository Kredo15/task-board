package board

import (
	"context"
	"strings"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type CreateBoardUseCase interface {
	Execute(ctx context.Context, cmd *CreateBoardRequest) (*CreateBoardResponse, error)
}

// createBoardHandler представляет обработчик команды создания доски
type createBoardUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewCreateBoardUseCase(r board.BoardRepository, g board.IDGenerator) CreateBoardUseCase {
	return &createBoardUseCase{
		repo: r,
		gen:  g,
	}
}

// Execute обрабатывает команду создания доски
func (h *createBoardUseCase) Execute(ctx context.Context, cmd *CreateBoardRequest) (*CreateBoardResponse, error) {
	// Преобразование запроса в доменную модель

	title := strings.TrimSpace(cmd.Title)
	desc := strings.TrimSpace(cmd.Description)
	if cmd.OwnerID == "" {
		return nil, board.ErrInvalidOwnerID
	}

	newBoard, err := board.NewBoard(
		h.gen,
		title,
		desc,
		cmd.OwnerID,
	)

	if err != nil {
		return nil, err
	}

	// Сохранение доски в репозитории
	if err := h.repo.Create(ctx, newBoard); err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &CreateBoardResponse{
		ID:        newBoard.ID(),
		Title:     newBoard.Title(),
		OwnerID:   newBoard.OwnerID(),
		CreatedAt: newBoard.CreatedAt(),
	}

	return response, nil
}
