package grpc

import (
	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	"github.com/Kredo15/task-board/services/board-service/pkg/logger"
	"github.com/Kredo15/task-board/services/board-service/pkg/validator"
)

type BoardHandler struct {
	boardv1.UnimplementedBoardServiceServer
	createUC CreateBoard
	getUC    GetBoard
	validate *validator.Validator
	log      logger.Logger
}

func NewBoardHandler(cUC CreateBoard, gUC GetBoard, v *validator.Validator, l logger.Logger) *BoardHandler {
	return &BoardHandler{
		createUC: cUC,
		getUC:    gUC,
		validate: v,
		log:      l,
	}
}
