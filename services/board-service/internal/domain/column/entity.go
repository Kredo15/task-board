package column

import (
	"time"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type Column struct {
	id        ColumnID
	boardID   board.BoardID
	title     Title
	rank      Rank
	createdAt time.Time
	updatedAt time.Time
}

func NewColumn(gen IDGenerator, boardRaw, titleRaw, rankRaw string) (*Column, error) {
	boardID, err := board.NewBoardID(boardRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Title
	title, err := NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Rank
	r, err := NewRank(rankRaw)
	if err != nil {
		return nil, err
	}

	board := &Column{
		id:        ColumnID(gen.Generate()),
		boardID:   boardID,
		title:     title,
		rank:      r,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	return board, nil
}

func RestoreColumn(idIn, boardIDIn, titleIn, rIn string, createdAt, updatedAt time.Time) (*Column, error) {
	id, err := NewColumnID(idIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Title
	title, err := NewTitle(titleIn)
	if err != nil {
		return nil, err
	}
	//
	bID, err := board.NewBoardID(boardIDIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Rank
	r, err := NewRank(rIn)
	if err != nil {
		return nil, err
	}

	return &Column{
		id:        id,
		boardID:   bID,
		title:     title,
		rank:      r,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (c *Column) ID() string { return string(c.id) }

func (c *Column) BoardID() string { return string(c.boardID) }

func (c *Column) Title() string { return string(c.title) }

func (c *Column) Rank() string { return string(c.rank) }

func (c *Column) CreatedAt() time.Time { return c.createdAt }

func (c *Column) UpdatedAt() time.Time { return c.updatedAt }

func (c *Column) Move(toRank Rank) {
	now := time.Now().UTC()

	c.rank = toRank
	c.updatedAt = now
}

func (c *Column) Update(title Title) {
	now := time.Now().UTC()

	c.title = title
	c.updatedAt = now
}

func (c *Column) Equals(other *Column) bool {
	if other == nil {
		return false
	}
	return c.id == other.id
}
