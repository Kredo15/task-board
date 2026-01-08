package board

import (
	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
)

// BoardHandler инкапсулирует gRPC клиент
type BoardHandler struct {
	boardClient boardv1.BoardServiceClient
}

func NewBoardHandler(client boardv1.BoardServiceClient) *BoardHandler {
	return &BoardHandler{boardClient: client}
}
