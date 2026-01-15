package task

import (
	"time"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

type Task struct {
	id          TaskID
	columnID    column.ColumnID
	title       Title
	description Description
	rank        Rank
	assigneeID  AssigneeID
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTask(gen IDGenerator, columnRaw, titleRaw, descRaw, rankRaw, assigneeRaw string) (*Task, error) {
	// Валидируем Title
	title, err := NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	//
	columnID, err := column.NewColumnID(columnRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Description
	desc, err := NewDescription(descRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Rank
	r, err := NewRank(rankRaw)
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
		columnID:    columnID,
		title:       title,
		description: desc,
		rank:        r,
		assigneeID:  assignee_id,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}

	return board, nil
}

func RestoreTask(id, colID, title, desc, r, assigneeID string, createdAt, updatedAt time.Time) *Task {
	return &Task{
		id:          TaskID(id),
		columnID:    column.ColumnID(colID),
		title:       Title(title),
		description: Description(desc),
		rank:        Rank(r),
		assigneeID:  AssigneeID(assigneeID),
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (t *Task) ID() string { return string(t.id) }

func (t *Task) ColumnID() string { return string(t.columnID) }

func (t *Task) Title() string { return string(t.title) }

func (t *Task) Description() string { return string(t.description) }

func (t *Task) Rank() string { return string(t.rank) }

func (t *Task) AssigneeID() string { return string(t.assigneeID) }

func (t *Task) CreatedAt() time.Time { return t.createdAt }

func (t *Task) UpdatedAt() time.Time { return t.updatedAt }

func (t *Task) UpdateTitle(newTitleRaw string) error {
	title, err := NewTitle(newTitleRaw)
	if err != nil {
		return err
	}
	t.title = title
	t.updatedAt = time.Now()
	return nil
}

func (t *Task) Move(toColumnId column.ColumnID, toRank Rank) {
	now := time.Now().UTC()

	t.columnID = toColumnId
	t.rank = toRank
	t.updatedAt = now
}

func (t *Task) UpdateDescription(newDescRaw string) error {
	desc, err := NewDescription(newDescRaw)
	if err != nil {
		return err
	}
	t.description = desc
	t.updatedAt = time.Now()
	return nil
}

func (t *Task) Equals(other *Task) bool {
	if other == nil {
		return false
	}
	return t.id == other.id
}
