package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	usecase "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

func mapDomainBoardToProto(b *usecase.CreateBoardResponse) *boardv1.Board {
	if b == nil {
		return nil
	}

	return &boardv1.Board{
		Id:          string(b.ID),
		Title:       b.Title,
		Description: b.Description,
		OwnerId:     b.OwnerID,
		CreatedAt:   timestamppb.New(b.CreatedAt),
		// Columns будут добавлены позже
	}
}
