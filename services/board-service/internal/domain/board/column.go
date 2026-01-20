package board

import (
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

func NewColumn(id ColumnID, boardID BoardID, title Title, rank Rank) *Column {

	column := &Column{
		id:        id,
		boardID:   boardID,
		title:     title,
		rank:      rank,
		tasks:     []*Task{},
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	return column
}

func RestoreColumn(
	idIn, bID, titleIn, rIn string,
	createdAt, updatedAt time.Time,
	tasks []*Task,
) (*Column, error) {
	// Валидируем ID
	id, err := ParseColumnID(idIn)
	if err != nil {
		return nil, err
	}
	// Валидируем BoardID
	boardID, err := ParseBoardID(bID)
	if err != nil {
		return nil, err
	}
	// Валидируем Title
	title, err := ParseTitle(titleIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Rank
	r, err := ParseRank(rIn)
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

func (c *Column) Update(title Title, eventID EventID) *DomainEvent {

	c.title = title
	c.updatedAt = time.Now()

	// Создаем Payload
	payload := ColumnUpdatedPayload{
		ColumnID: c.ID(),
		BoardID:  c.BoardID(),
		NewTitle: string(c.title),
	}

	return &DomainEvent{
		ID:          eventID,
		Type:        EventTypeColumnUpdated,
		AggregateID: c.BoardID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}
}

// Move изменяет позицию (Rank) колонки
func (c *Column) Move(newRank Rank, boardID BoardID, eventID EventID) *DomainEvent {
	c.rank = newRank
	c.updatedAt = time.Now()

	// Создаем Payload
	payload := ColumnMovedPayload{
		ColumnID: c.ID(),
		BoardID:  c.BoardID(),
		NewRank:  c.Rank(),
	}

	// Возвращаем событие
	return &DomainEvent{
		ID:          eventID, // ID самого события
		Type:        EventTypeColumnMoved,
		AggregateID: c.BoardID(),
		Payload:     payload,
		OccurredAt:  c.updatedAt,
	}
}

func (c *Column) SetTasks(tasks []*Task) {
	c.tasks = tasks
}

func (c *Column) Delete(eventID EventID) *DomainEvent {
	// Создаем Payload
	payload := ColumnDeletedPayload{
		ColumnID: c.ID(),
		BoardID:  c.BoardID(),
	}

	// Создаем событие
	return &DomainEvent{
		ID:          eventID,
		Type:        EventTypeColumnDeleted,
		AggregateID: c.BoardID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}
}
