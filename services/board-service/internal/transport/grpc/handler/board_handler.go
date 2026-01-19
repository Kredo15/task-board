package grpc

import (
	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"
)

type BoardHandler struct {
	boardv1.UnimplementedBoardServiceServer
	createBoardUC CreateBoard
	getBoardUC    GetBoard
	validate      *validator.Validator
	log           logger.Logger
}

func NewBoardHandler(createBoardUC CreateBoard, getBoardUC GetBoard, validate *validator.Validator, log logger.Logger) *BoardHandler {
	return &BoardHandler{
		createBoardUC: createBoardUC,
		getBoardUC:    getBoardUC,
		validate:      validate,
		log:           log,
	}
}
