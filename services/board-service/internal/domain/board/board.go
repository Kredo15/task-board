package board

import (
	"time"
)

const MaxColumnsPerBoard = 10
const MaxTasksPerColumn = 50

type Board struct {
	id          BoardID
	title       Title
	description Description
	ownerID     OwnerID
	columns     []*Column
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBoard(id BoardID, title Title, desc Description, ownerID OwnerID, eventID EventID) (*Board, *DomainEvent) {

	board := &Board{
		id:          id,
		title:       title,
		description: desc,
		ownerID:     ownerID,
		columns:     []*Column{},
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}

	// Создаем Payload
	payload := BoardCreatedPayload{
		BoardID: board.ID(),
		Title:   board.Title(),
		OwnerID: board.OwnerID(),
	}

	event := &DomainEvent{
		ID:          eventID,
		Type:        EventTypeBoardCreated,
		AggregateID: board.ID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}

	return board, event
}

func RestoreBoard(
	id, title, desc, ownerID string,
	createdAt, updatedAt time.Time,
	columns []*Column,
) *Board {
	return &Board{
		id:          BoardID(id),
		title:       Title(title),
		description: Description(desc),
		ownerID:     OwnerID(ownerID),
		columns:     columns,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (b *Board) ID() string           { return string(b.id) }
func (b *Board) Title() string        { return string(b.title) }
func (b *Board) Description() string  { return string(b.description) }
func (b *Board) OwnerID() string      { return string(b.ownerID) }
func (b *Board) Columns() []*Column   { return b.columns }
func (b *Board) CreatedAt() time.Time { return b.createdAt }
func (b *Board) UpdatedAt() time.Time { return b.updatedAt }

func (b *Board) Update(title *Title, desc *Description, eventID EventID) (*DomainEvent, error) {
	var payloadTitle *string
	var payloadDesc *string

	// Обновляем заголовок, если он передан
	if title != nil {
		b.title = *title
		s := string(*title)
		payloadTitle = &s
	}

	// Обновляем описание, если оно передано
	if desc != nil {
		b.description = *desc
		s := string(*desc)
		payloadDesc = &s
	}

	b.updatedAt = time.Now()

	// Создаем Payload
	payload := BoardUpdatedPayload{
		BoardID:     b.ID(),
		Title:       payloadTitle,
		Description: payloadDesc,
	}

	return &DomainEvent{
		ID:          eventID,
		Type:        EventTypeBoardUpdated,
		AggregateID: b.ID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}, nil
}

func (b *Board) AddColumn(columnID ColumnID, currentColumnCount int, title Title, rank Rank, eventID EventID) (*Column, *DomainEvent, error) {
	// проверка лимита колонок
	if currentColumnCount >= MaxColumnsPerBoard {
		return nil, nil, ErrLimitColumnReached
	}
	// Создаем объект колонки
	column := NewColumn(
		columnID,
		b.id,
		title,
		rank,
	)

	b.columns = append(b.columns, column)
	b.updatedAt = time.Now()

	// Создаем Payload
	payload := ColumnCreatedPayload{
		ColumnID: column.ID(),
		BoardID:  b.ID(),
		Title:    column.Title(),
		Rank:     column.Rank(),
	}

	// Создаем событие
	event := &DomainEvent{
		ID:          eventID,
		Type:        EventTypeColumnCreated,
		AggregateID: b.ID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}
	return column, event, nil
}

func (b *Board) AddTask(
	taskID TaskID,
	currentTaskCount int,
	columnID ColumnID,
	title Title,
	description Description,
	rank Rank,
	assigneeId AssigneeID,
	eventID EventID,
) (*Task, *DomainEvent, error) {

	// проверка лимита задач в колонке
	if currentTaskCount >= MaxTasksPerColumn {
		return nil, nil, ErrLimitTaskReached
	}
	// Создаем объект задачи
	task := NewTask(
		taskID,
		b.id,
		columnID,
		title,
		description,
		rank,
		assigneeId,
	)

	// Создаем Payload
	payload := TaskCreatedPayload{
		TaskID:      task.ID(),
		ColumnID:    string(columnID),
		BoardID:     b.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		Rank:        task.Rank(),
	}

	// Создаем событие
	event := &DomainEvent{
		ID:          eventID,
		Type:        EventTypeTaskCreated,
		AggregateID: b.ID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}
	return task, event, nil
}

func (b *Board) SetColumns(cols []*Column) {
	b.columns = cols
}

func (b *Board) Delete(eventID EventID) *DomainEvent {
	// Создаем Payload
	payload := BoardDeletedPayload{
		BoardID: b.ID(),
	}

	return &DomainEvent{
		ID:          eventID,
		Type:        EventTypeBoardDeleted,
		AggregateID: b.ID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}
}
