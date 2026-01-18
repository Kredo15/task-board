package board

import (
	"encoding/json"
	"time"
)

type Board struct {
	id          BoardID
	title       Title
	description Description
	ownerID     OwnerID
	columns     []*Column
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBoard(gen IDGenerator, titleRaw, descRaw, ownerRaw string) (*Board, *DomainEvent, error) {
	// Валидируем Title
	title, err := NewTitle(titleRaw)
	if err != nil {
		return nil, nil, err
	}
	// Валидируем Description
	desc, err := NewDescription(descRaw)
	if err != nil {
		return nil, nil, err
	}
	// Валидируем OwnerID
	owner_id, err := NewOwnerID(ownerRaw)
	if err != nil {
		return nil, nil, err
	}

	board := &Board{
		id:          BoardID(gen.Generate()),
		title:       title,
		description: desc,
		ownerID:     owner_id,
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

	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	event := &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeBoardCreated,
		AggregateID: board.ID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}

	return board, event, nil
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

func (b *Board) Update(gen IDGenerator, newTitleRaw, newDescRaw string) (*DomainEvent, error) {
	title, err := NewTitle(newTitleRaw)
	if err != nil {
		return nil, err
	}

	description, err := NewDescription(newDescRaw)
	if err != nil {
		return nil, err
	}

	b.title = title
	b.description = description
	b.updatedAt = time.Now()

	// Создаем Payload
	payload := BoardUpdatedPayload{
		BoardID:     b.ID(),
		Title:       &newTitleRaw,
		Description: &newDescRaw,
	}

	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeBoardUpdated,
		AggregateID: b.ID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}, nil
}

func (b *Board) AddColumn(gen IDGenerator, titleStr, lastRankstr string) (*Column, *DomainEvent, error) {
	title, err := NewTitle(titleStr)
	if err != nil {
		return nil, nil, err
	}
	rank, err := NewRank(lastRankstr)
	if err != nil {
		return nil, nil, err
	}

	// Создаем объект колонки
	column, err := RestoreColumn(
		gen.Generate(),
		b.ID(),
		string(title),
		string(rank),
		time.Now(),
		time.Now(),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	// Создаем Payload
	payload := ColumnCreatedPayload{
		ColumnID: column.ID(),
		BoardID:  b.ID(),
		Title:    column.Title(),
		Rank:     column.Rank(),
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	// Создаем событие
	event := &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeColumnCreated,
		AggregateID: b.ID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}
	return column, event, nil
}

func (b *Board) AddTask(
	gen IDGenerator,
	currentTaskCount int,
	colID, titleStr, desc, lastRankstr, assigneeId string,
) (*Task, *DomainEvent, error) {

	// проверка лимита задач в колонке
	if currentTaskCount >= 50 {
		return nil, nil, ErrLimitTaskReached
	}
	columnID, err := NewColumnID(colID)
	if err != nil {
		return nil, nil, err
	}

	title, err := NewTitle(titleStr)
	if err != nil {
		return nil, nil, err
	}
	lastRank, err := NewRank(lastRankstr)
	if err != nil {
		return nil, nil, err
	}
	description, err := NewDescription(desc)
	if err != nil {
		return nil, nil, err
	}
	// Создаем объект задачи
	task, err := RestoreTask(
		gen.Generate(),
		b.ID(),
		string(columnID),
		string(title),
		string(description),
		string(lastRank),
		assigneeId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, nil, err
	}

	// Создаем Payload
	payload := TaskCreatedPayload{
		TaskID:      task.ID(),
		ColumnID:    string(columnID),
		BoardID:     b.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		Rank:        task.Rank(),
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	// Создаем событие
	event := &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeTaskCreated,
		AggregateID: b.ID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}
	return task, event, nil
}

func (b *Board) SetColumns(cols []*Column) {
	b.columns = cols
}

func (b *Board) Delete(gen IDGenerator) (*DomainEvent, error) {
	// Создаем Payload
	payload := BoardDeletedPayload{
		BoardID: b.ID(),
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeBoardDeleted,
		AggregateID: b.ID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}, nil
}
