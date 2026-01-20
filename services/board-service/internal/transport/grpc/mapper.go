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

func mapDTOBoardViewToProto(b *usecase.BoardResponse) *boardv1.BoardView {
	if b == nil {
		return nil
	}

	columns := make([]*boardv1.ColumnView, 0, len(b.Columns))

	for _, col := range b.Columns {
		// Сначала маппим колонку
		column := &boardv1.Column{
			Id:        col.ID,
			Title:     col.Title,
			Rank:      col.Rank,
			CreatedAt: timestamppb.New(col.CreatedAt),
			UpdatedAt: timestamppb.New(col.UpdatedAt),
			Tasks:     make([]*boardv1.Task, 0, len(col.Tasks)),
		}
		// Затем маппим задачи внутри колонки
		for _, task := range col.Tasks {
			taskProto := &boardv1.Task{
				Id:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				Rank:        task.Rank,
				CreatedAt:   timestamppb.New(task.CreatedAt),
				UpdatedAt:   timestamppb.New(task.UpdatedAt),
			}
			column.Tasks = append(column.Tasks, taskProto)
		}
		// Создаем ColumnView и добавляем в слайс
		columnView := &boardv1.ColumnView{
			Column: column,
		}
		columns = append(columns, columnView)
	}

	board := &boardv1.Board{
		Id:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		OwnerId:     b.OwnerID,
		CreatedAt:   timestamppb.New(b.CreatedAt),
		UpdatedAt:   timestamppb.New(b.UpdatedAt),
	}

	boardView := &boardv1.BoardView{
		Board:   board,
		Columns: columns,
	}
	return boardView
}
