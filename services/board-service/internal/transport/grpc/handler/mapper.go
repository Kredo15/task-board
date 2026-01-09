package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	boardv1 "github.com/Kredo15/task-board/protos/gen/go/board/v1"
	usecase "github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

func mapDTOBoardToProto(b *usecase.BoardResponse) *boardv1.Board {
	if b == nil {
		return nil
	}

	return &boardv1.Board{
		Id:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		OwnerId:     b.OwnerID,
		CreatedAt:   timestamppb.New(b.CreatedAt),
	}
}

func mapDTOBoardViewToProto(b *usecase.BoardResponse) *boardv1.Board {
	if b == nil {
		return nil
	}

	board := &boardv1.Board{
		Id:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		OwnerId:     b.OwnerID,
		CreatedAt:   timestamppb.New(b.CreatedAt),
		UpdatedAt:   timestamppb.New(b.CreatedAt),
	}
	columns := make([]*boardv1.ColumnView, 0)
	boardView := &boardv1.BoardView{
		Board:   board,
		Columns: make([]*boardv1.ColumnView, 0), // Columns добавлю позже
	}
	return boardView
}
