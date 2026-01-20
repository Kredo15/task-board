package board

import (
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
	taskID TaskID,
	boardID BoardID,
	colID ColumnID,
	title Title,
	desc Description,
	rank Rank,
	assignee_id AssigneeID,
) *Task {

	task := &Task{
		id:          taskID,
		boardID:     boardID,
		columnID:    colID,
		title:       title,
		description: desc,
		rank:        rank,
		assigneeID:  assignee_id,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}

	return task
}

func RestoreTask(
	idIn, boardIDIn, colIDIn, titleIn, descIn, rIn, assigneeIDIn string,
	createdAt, updatedAt time.Time,
) (*Task, error) {
	id, err := ParseTaskID(idIn)
	if err != nil {
		return nil, err
	}
	boardID, err := ParseBoardID(boardIDIn)
	if err != nil {
		return nil, err
	}
	colID, err := ParseColumnID(colIDIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Title
	title, err := ParseTitle(titleIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Description
	desc, err := ParseDescription(descIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Rank
	r, err := ParseRank(rIn)
	if err != nil {
		return nil, err
	}
	// Валидируем AssigneeID
	assignee_id, err := ParseAssigneeID(assigneeIDIn)
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

func (t *Task) Update(title *Title, desc *Description, assigneeID *AssigneeID, eventID EventID) *DomainEvent {
	var payloadTitle *string
	var payloadDesc *string
	var payloadAssigneeID *string
	// Обновляем заголовок, если он передан
	if title != nil {
		t.title = *title
		s := string(*title)
		payloadTitle = &s
	}

	// Обновляем описание, если оно передано
	if desc != nil {
		t.description = *desc
		s := string(*desc)
		payloadDesc = &s
	}
	// Обновляем исполнителя, если он передан
	if assigneeID != nil {
		t.assigneeID = *assigneeID
		s := string(*assigneeID)
		payloadAssigneeID = &s
	}

	t.updatedAt = time.Now()

	// Создаем Payload
	payload := TaskUpdatedPayload{
		TaskID:      t.ID(),
		ColumnID:    t.ColumnID(),
		BoardID:     t.BoardID(),
		Title:       payloadTitle,
		Description: payloadDesc,
		AssigneeID:  payloadAssigneeID,
	}

	// Возвращаем событие
	return &DomainEvent{
		ID:          eventID, // ID самого события
		Type:        EventTypeTaskUpdated,
		AggregateID: t.BoardID(),
		Payload:     payload,
		OccurredAt:  t.updatedAt,
	}
}

func (t *Task) Move(toColumnID ColumnID, newRank Rank, eventID EventID) *DomainEvent {
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

	// Возвращаем событие
	return &DomainEvent{
		ID:          eventID, // ID самого события
		Type:        EventTypeTaskMoved,
		AggregateID: t.BoardID(),
		Payload:     payload,
		OccurredAt:  t.updatedAt,
	}
}

func (t *Task) Delete(eventID EventID) *DomainEvent {
	// Создаем Payload
	payload := TaskDeletedPayload{
		TaskID:   t.ID(),
		ColumnID: t.ColumnID(),
		BoardID:  t.BoardID(),
	}

	// Создаем событие
	return &DomainEvent{
		ID:          eventID,
		Type:        EventTypeTaskDeleted,
		AggregateID: t.BoardID(),
		Payload:     payload,
		OccurredAt:  time.Now(),
	}
}
