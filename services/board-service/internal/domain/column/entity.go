package column

import (
	"time"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type Column struct {
	id        ColumnID
	boardID   board.BoardID
	title     Title
	position  Position
	createdAt time.Time
	updatedAt time.Time
}

func NewBoard(gen IDGenerator, titleRaw string, posRaw int) (*Column, error) {
	// Валидируем Title
	title, err := NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Position
	pos, err := NewPosition(posRaw)
	if err != nil {
		return nil, err
	}

	board := &Column{
		id:        ColumnID(gen.Generate()),
		title:     title,
		position:  pos,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	return board, nil
}

func (c *Column) ID() string { return string(c.id) }

func (c *Column) Title() string { return string(c.title) }

func (c *Column) CreatedAt() time.Time { return c.createdAt }

func (c *Column) UpdatedAt() time.Time { return c.updatedAt }
