package task

import "time"

type Task struct {
	id          TaskID
	columnID    string
	title       Title
	description Description
	position    Position
	assigneeID  AssigneeID
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBoard(gen IDGenerator, titleRaw, descRaw, assigneeRaw string, posRaw int) (*Task, error) {
	// Валидируем Title
	title, err := NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Description
	desc, err := NewDescription(descRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Position
	pos, err := NewPosition(posRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем AssigneeID
	assignee_id, err := NewAssigneeID(assigneeRaw)
	if err != nil {
		return nil, err
	}

	board := &Task{
		id:          TaskID(gen.Generate()),
		title:       title,
		description: desc,
		position:    pos,
		assigneeID:  assignee_id,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}

	return board, nil
}

func RestoreTask(id, title, desc, assigneeID string, pos int, createdAt, updatedAt time.Time) *Task {
	return &Task{
		id:          TaskID(id),
		title:       Title(title),
		description: Description(desc),
		position:    Position(pos),
		assigneeID:  AssigneeID(assigneeID),
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (b *Task) ID() string { return string(b.id) }

func (b *Task) Title() string { return string(b.title) }

func (b *Task) Description() string { return string(b.description) }

func (b *Task) AssigneeID() string { return string(b.assigneeID) }

func (b *Task) CreatedAt() time.Time { return b.createdAt }

func (b *Task) UpdatedAt() time.Time { return b.updatedAt }

func (b *Task) UpdateTitle(newTitleRaw string) error {
	title, err := NewTitle(newTitleRaw)
	if err != nil {
		return err
	}
	b.title = title
	b.updatedAt = time.Now()
	return nil
}

func (b *Task) UpdateDescription(newDescRaw string) error {
	desc, err := NewDescription(newDescRaw)
	if err != nil {
		return err
	}
	b.description = desc
	b.updatedAt = time.Now()
	return nil
}

func (b *Task) Equals(other *Task) bool {
	if other == nil {
		return false
	}
	return b.id == other.id
}
