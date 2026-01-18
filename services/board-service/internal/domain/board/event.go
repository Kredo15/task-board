package board

import (
	"time"
)

// Константы типов событий
const (
	EventTypeBoardCreated  = "board.created"
	EventTypeBoardUpdated  = "board.updated"
	EventTypeBoardDeleted  = "board.deleted"
	EventTypeColumnCreated = "column.created"
	EventTypeColumnUpdated = "column.updated"
	EventTypeColumnMoved   = "column.moved"
	EventTypeColumnDeleted = "column.deleted"
	EventTypeTaskCreated   = "task.created"
	EventTypeTaskMoved     = "task.moved"
	EventTypeTaskUpdated   = "task.updated"
	EventTypeTaskDeleted   = "task.deleted"
)

// DomainEvent — общая структура для хранения в Outbox
type DomainEvent struct {
	ID          string
	Type        string
	AggregateID string
	Payload     []byte
	OccurredAt  time.Time
}

// --- Board Events ---

type BoardCreatedPayload struct {
	BoardID string `json:"board_id"`
	Title   string `json:"title"`
	OwnerID string `json:"owner_id"`
}

type BoardUpdatedPayload struct {
	BoardID     string  `json:"board_id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type BoardDeletedPayload struct {
	BoardID string `json:"board_id"`
}

// --- Column Events ---

type ColumnCreatedPayload struct {
	ColumnID string `json:"column_id"`
	BoardID  string `json:"board_id"`
	Title    string `json:"title"`
	Rank     string `json:"rank"`
}

type ColumnUpdatedPayload struct {
	ColumnID string `json:"column_id"`
	BoardID  string `json:"board_id"`
	NewTitle string `json:"new_title"`
}

type ColumnMovedPayload struct {
	ColumnID string `json:"column_id"`
	BoardID  string `json:"board_id"`
	NewRank  string `json:"new_rank"`
}

type ColumnDeletedPayload struct {
	BoardID  string `json:"board_id"`
	ColumnID string `json:"column_id"`
}

// --- Task Events ---

type TaskCreatedPayload struct {
	TaskID      string `json:"task_id"`
	ColumnID    string `json:"column_id"`
	BoardID     string `json:"board_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Rank        string `json:"rank"`
	AssigneeID  string `json:"assignee_id"`
}

type TaskUpdatedPayload struct {
	TaskID      string  `json:"task_id"`
	ColumnID    string  `json:"column_id"`
	BoardID     string  `json:"board_id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	AssigneeID  *string `json:"assignee_id,omitempty"`
}

type TaskMovedPayload struct {
	TaskID       string `json:"task_id"`
	ColumnID     string `json:"column_id"`
	BoardID      string `json:"board_id"`
	FromColumnID string `json:"from_column_id"`
	ToColumnID   string `json:"to_column_id"`
	NewRank      string `json:"new_rank"`
}

type TaskDeletedPayload struct {
	TaskID   string `json:"task_id"`
	ColumnID string `json:"column_id"`
	BoardID  string `json:"board_id"`
}
