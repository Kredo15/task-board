package board

import (
	"encoding/json"
	"time"
)

type Column struct {
	id        ColumnID
	boardID   BoardID
	title     Title
	rank      Rank
	tasks     []*Task
	createdAt time.Time
	updatedAt time.Time
}

func NewColumn(gen IDGenerator, boardRaw, titleRaw, rankRaw string) (*Column, error) {
	boardID, err := NewBoardID(boardRaw)
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

func RestoreColumn(
	idIn, bID, titleIn, rIn string,
	createdAt, updatedAt time.Time,
	tasks []*Task,
) (*Column, error) {
	// Валидируем ID
	id, err := NewColumnID(idIn)
	if err != nil {
		return nil, err
	}
	// Валидируем BoardID
	boardID, err := NewBoardID(bID)
	if err != nil {
		return nil, err
	}
	// Валидируем Title
	title, err := NewTitle(titleIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Rank
	r, err := NewRank(rIn)
	if err != nil {
		return nil, err
	}

	return &Column{
		id:        ColumnID(id),
		boardID:   boardID,
		title:     Title(title),
		rank:      Rank(r),
		tasks:     tasks,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (c *Column) ID() string           { return string(c.id) }
func (c *Column) BoardID() string      { return string(c.boardID) }
func (c *Column) Title() string        { return string(c.title) }
func (c *Column) Rank() string         { return string(c.rank) }
func (c *Column) Tasks() []*Task       { return c.tasks }
func (c *Column) CreatedAt() time.Time { return c.createdAt }
func (c *Column) UpdatedAt() time.Time { return c.updatedAt }

func (c *Column) Update(gen IDGenerator, newTitleRaw string) (*DomainEvent, error) {
	newTitle, err := NewTitle(newTitleRaw)
	if err != nil {
		return nil, err
	}
	c.title = newTitle
	c.updatedAt = time.Now()

	// Создаем Payload
	payload := ColumnUpdatedPayload{
		ColumnID: c.ID(),
		BoardID:  c.BoardID(),
		NewTitle: newTitleRaw,
	}

	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeColumnUpdated,
		AggregateID: c.BoardID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}, nil
}

// Move изменяет позицию (Rank) колонки
func (c *Column) Move(gen IDGenerator, newRank Rank, boardID BoardID) (*DomainEvent, error) {
	c.rank = newRank
	c.updatedAt = time.Now()

	// Создаем Payload
	payload := ColumnMovedPayload{
		ColumnID: c.ID(),
		BoardID:  c.BoardID(),
		NewRank:  c.Rank(),
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// Возвращаем событие
	return &DomainEvent{
		ID:          gen.Generate(), // ID самого события
		Type:        EventTypeColumnMoved,
		AggregateID: c.BoardID(),
		Payload:     data,
		OccurredAt:  c.updatedAt,
	}, nil
}

func (c *Column) SetTasks(tasks []*Task) {
	c.tasks = tasks
}

func (c *Column) Delete(gen IDGenerator) (*DomainEvent, error) {
	// Создаем Payload
	payload := ColumnDeletedPayload{
		ColumnID: c.ID(),
		BoardID:  c.BoardID(),
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// Создаем событие
	return &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeColumnDeleted,
		AggregateID: c.BoardID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}, nil
}
