package task

import "time"

type CreateTaskRequest struct {
	ColumnID    string `json:"column_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Position    int32  `json:"position" validate:"required"`
	AssigneeID  string `json:"assignee_id" validate:"required"`
}

type MoveTaskRequest struct {
	ID       string `json:"id" validate:"required"`
	ColumnID string `json:"column_id" validate:"required"`
	Position int32  `json:"position" validate:"required"`
}

type DeleteTaskRequest struct {
	ID string `json:"id" validate:"required"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	ColumnID    string    `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int32     `json:"position"`
	AssigneeID  string    `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeleteTaskResponse struct{}
