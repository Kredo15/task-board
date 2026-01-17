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

func RestoreTask(idIn, colIDIn, titleIn, descIn, rIn, assigneeIDIn string, createdAt, updatedAt time.Time) (*Task, error) {
	id, err := NewTaskID(idIn)
	if err != nil {
		return nil, err
	}
	// Валидируем Title
	title, err := NewTitle(titleIn)
	if err != nil {
		return nil, err
	}
	//
	colID, err := column.NewColumnID(colIDIn)
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
		columnID:    colID,
		title:       title,
		description: desc,
		rank:        r,
		assigneeID:  assignee_id,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

func (t *Task) ID() string { return string(t.id) }

func (t *Task) ColumnID() string { return string(t.columnID) }

func (t *Task) Title() string { return string(t.title) }

func (t *Task) Description() string { return string(t.description) }

func (t *Task) Rank() string { return string(t.rank) }

func (t *Task) AssigneeID() string { return string(t.assigneeID) }

func (t *Task) CreatedAt() time.Time { return t.createdAt }

func (t *Task) UpdatedAt() time.Time { return t.updatedAt }

func (t *Task) Move(toColumnId column.ColumnID, toRank Rank) {
	now := time.Now().UTC()

	t.columnID = toColumnId
	t.rank = toRank
	t.updatedAt = now
}

func (t *Task) Update(title Title, desc Description, assigneeID AssigneeID) {
	now := time.Now().UTC()

	t.title = title
	t.description = desc
	t.assigneeID = assigneeID
	t.updatedAt = now
}

func (t *Task) Equals(other *Task) bool {
	if other == nil {
		return false
	}
	return t.id == other.id
}
