package board

import (
	"encoding/json"
	"time"
)

type Task struct {
	id          TaskID
	boardID     BoardID
	columnID    ColumnID
	title       Title
	description Description
	rank        Rank
	assigneeID  AssigneeID
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTask(
	gen IDGenerator, boardIDIn, colIDIn, titleRaw, descRaw, rankRaw, assigneeIDRaw string,
) (*Task, *DomainEvent, error) {
	boardID, err := NewBoardID(boardIDIn)
	if err != nil {
		return nil, nil, err
	}
	colID, err := NewColumnID(colIDIn)
	if err != nil {
		return nil, nil, err
	}
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
	// Валидируем Rank
	r, err := NewRank(rankRaw)
	if err != nil {
		return nil, nil, err
	}
	// Валидируем AssigneeID
	assignee_id, err := NewAssigneeID(assigneeIDRaw)
	if err != nil {
		return nil, nil, err
	}
	task := &Task{
		id:          TaskID(gen.Generate()),
		boardID:     boardID,
		columnID:    colID,
		title:       title,
		description: desc,
		rank:        r,
		assigneeID:  assignee_id,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
	// Создаем Payload
	payload := TaskCreatedPayload{
		TaskID:     task.ID(),
		ColumnID:   task.ColumnID(),
		Title:      task.Title(),
		Rank:       task.Rank(),
		AssigneeID: task.AssigneeID(),
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	// Возвращаем задачу и событие
	event := &DomainEvent{
		ID:          gen.Generate(), // ID самого события
		Type:        EventTypeTaskCreated,
		AggregateID: task.BoardID(),
		Payload:     data,
		OccurredAt:  task.CreatedAt(),
	}
	return task, event, nil
}

func RestoreTask(
	idIn, boardIDIn, colIDIn, titleIn, descIn, rIn, assigneeIDIn string,
	createdAt, updatedAt time.Time,
) (*Task, error) {
	id, err := NewTaskID(idIn)
	if err != nil {
		return nil, err
	}
	boardID, err := NewBoardID(boardIDIn)
	if err != nil {
		return nil, err
	}
	colID, err := NewColumnID(colIDIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Title
	title, err := NewTitle(titleIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Description
	desc, err := NewDescription(descIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Rank
	r, err := NewRank(rIn)
	if err != nil {
		return nil, err
	}
	// Валидируем AssigneeID
	assignee_id, err := NewAssigneeID(assigneeIDIn)
	if err != nil {
		return nil, err
	}

	return &Task{
		id:          id,
		boardID:     boardID,
		columnID:    colID,
		title:       title,
		description: desc,
		rank:        r,
		assigneeID:  assignee_id,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

func (t *Task) ID() string           { return string(t.id) }
func (t *Task) BoardID() string      { return string(t.boardID) }
func (t *Task) ColumnID() string     { return string(t.columnID) }
func (t *Task) Title() string        { return string(t.title) }
func (t *Task) Description() string  { return string(t.description) }
func (t *Task) Rank() string         { return string(t.rank) }
func (t *Task) AssigneeID() string   { return string(t.assigneeID) }
func (t *Task) CreatedAt() time.Time { return t.createdAt }
func (t *Task) UpdatedAt() time.Time { return t.updatedAt }

func (t *Task) Update(
	gen IDGenerator,
	newTitleRaw, newDescRaw, newAssigneeIDRaw string,
) (*DomainEvent, error) {
	newTitle, err := NewTitle(newTitleRaw)
	if err != nil {
		return nil, err
	}
	newDesc, err := NewDescription(newDescRaw)
	if err != nil {
		return nil, err
	}
	newAssigneeID, err := NewAssigneeID(newAssigneeIDRaw)
	if err != nil {
		return nil, err
	}

	t.title = newTitle
	t.description = newDesc
	t.assigneeID = newAssigneeID
	t.updatedAt = time.Now()

	// Создаем Payload
	payload := TaskUpdatedPayload{
		TaskID:      t.ID(),
		ColumnID:    t.ColumnID(),
		BoardID:     t.BoardID(),
		Title:       &newTitleRaw,
		Description: &newDescRaw,
		AssigneeID:  &newAssigneeIDRaw,
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// Возвращаем событие
	return &DomainEvent{
		ID:          gen.Generate(), // ID самого события
		Type:        EventTypeTaskUpdated,
		AggregateID: t.BoardID(),
		Payload:     data,
		OccurredAt:  t.updatedAt,
	}, nil
}

func (t *Task) Move(gen IDGenerator, toColumnID ColumnID, newRank Rank) (*DomainEvent, error) {
	oldColID := t.columnID
	t.columnID = toColumnID
	t.rank = newRank
	t.updatedAt = time.Now()

	// Создаем Payload
	payload := TaskMovedPayload{
		TaskID:       t.ID(),
		ColumnID:     string(t.columnID),
		BoardID:      t.BoardID(),
		FromColumnID: string(oldColID),
		ToColumnID:   string(toColumnID),
		NewRank:      t.Rank(),
	}

	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Возвращаем событие
	return &DomainEvent{
		ID:          gen.Generate(), // ID самого события
		Type:        EventTypeTaskMoved,
		AggregateID: t.BoardID(),
		Payload:     data,
		OccurredAt:  t.updatedAt,
	}, nil
}

func (t *Task) Delete(gen IDGenerator) (*DomainEvent, error) {
	// Создаем Payload
	payload := TaskDeletedPayload{
		TaskID:   t.ID(),
		ColumnID: t.ColumnID(),
		BoardID:  t.BoardID(),
	}
	// Сериализуем Payload
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// Создаем событие
	return &DomainEvent{
		ID:          gen.Generate(),
		Type:        EventTypeTaskDeleted,
		AggregateID: t.BoardID(),
		Payload:     data,
		OccurredAt:  time.Now(),
	}, nil
}
