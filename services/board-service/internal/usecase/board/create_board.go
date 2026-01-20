package board

import (
	"context"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

// createBoardHandler представляет обработчик команды создания доски
type CreateBoardUseCase struct {
	repo board.BoardRepository
	gen  board.IDGenerator
}

// NewCreateBoardHandler создает новый экземпляр обработчика команды создания доски
func NewCreateBoardUseCase(r board.BoardRepository, g board.IDGenerator) *CreateBoardUseCase {
	return &CreateBoardUseCase{
		repo: r,
		gen:  g,
	}
}

// Execute обрабатывает команду создания доски
func (uc *CreateBoardUseCase) Execute(ctx context.Context, req *CreateBoardRequest) (*BoardResponse, error) {
	title, err := board.ParseTitle(req.Title)
	if err != nil {
		return nil, err // Возвращаем ошибку домена (например, ErrInvalidBoardTitleEmpty)
	}

	desc, err := board.ParseDescription(req.Description)
	if err != nil {
		return nil, err
	}

	owner, err := board.ParseOwnerID(req.OwnerID)
	if err != nil {
		return nil, err
	}

	// 2. Только когда всё валидно, генерируем ID и создаем домен
	boardID := board.BoardID(uc.gen.Generate())
	eventID := board.EventID(uc.gen.Generate())

	newBoard, event := board.NewBoard(boardID, title, desc, owner, eventID)

	// Сохраненяем доску в репозитории и сохраняем событие
	if err := uc.repo.Create(ctx, newBoard, event); err != nil {
		return nil, err
	}

	// Возвращаем успешный ответ
	response := &BoardResponse{
		ID:          newBoard.ID(),
		Title:       newBoard.Title(),
		Description: newBoard.Description(),
		OwnerID:     newBoard.OwnerID(),
		CreatedAt:   newBoard.CreatedAt(),
		UpdatedAt:   newBoard.UpdatedAt(),
	}

	return response, nil
}
