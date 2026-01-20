package grpc

import (
	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"
)

type BoardUseCases struct {
	Create *board.CreateBoardUseCase
	Update *board.UpdateBoardUseCase
	Get    *board.GetBoardUseCase
	Delete *board.DeleteBoardUseCase
}

type ColumnUseCases struct {
	Create *column.AddColumnUseCase
	Update *column.UpdateColumnUseCase
	Move   *column.MoveColumnUseCase
	Delete *column.DeleteColumnUseCase
}

type TaskUseCases struct {
	Create *task.AddTaskUseCase
	Update *task.UpdateTaskUseCase
	Move   *task.MoveTaskUseCase
	Delete *task.DeleteTaskUseCase
}

// Deps — структура для внедрения зависимостей
type Deps struct {
	Board  BoardUseCases
	Column ColumnUseCases
	Task   TaskUseCases
}

type Handler struct {
	boardv1.UnimplementedBoardServiceServer
	board    BoardUseCases
	column   ColumnUseCases
	task     TaskUseCases
	validate *validator.Validator
	log      logger.Logger
}

func NewHandler(d Deps, validate *validator.Validator, log logger.Logger) *Handler {
	return &Handler{
		board:    d.Board,
		column:   d.Column,
		task:     d.Task,
		validate: validate,
		log:      log,
	}
}
