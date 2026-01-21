package grpc

import (
	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
	"github.com/Kredo15/task-board/services/board-service/internal/usecase/column"
	"github.com/Kredo15/task-board/services/board-service/internal/usecase/task"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
)

type BoardUseCases struct {
	Create board.CreateBoard
	Update board.UpdateBoard
	Get    board.GetBoard
	Delete board.DeleteBoard
}

type ColumnUseCases struct {
	Create column.AddColumn
	Update column.UpdateColumn
	Move   column.MoveColumn
	Delete column.DeleteColumn
}

type TaskUseCases struct {
	Create task.AddTask
	Update task.UpdateTask
	Move   task.MoveTask
	Delete task.DeleteTask
}

// Deps — структура для внедрения зависимостей
type Deps struct {
	Board  BoardUseCases
	Column ColumnUseCases
	Task   TaskUseCases
}

type Handler struct {
	boardv1.UnimplementedBoardServiceServer
	board  BoardUseCases
	column ColumnUseCases
	task   TaskUseCases
	log    logger.Logger
}

func NewHandler(d Deps, log logger.Logger) *Handler {
	return &Handler{
		board:  d.Board,
		column: d.Column,
		task:   d.Task,
		log:    log,
	}
}
