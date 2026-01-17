package task

import "time"

type CreateTaskRequest struct {
	ColumnID    string `json:"column_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Rank        string `json:"rank" validate:"required"`
	AssigneeID  string `json:"assignee_id" validate:"required"`
	BoardID     string `json:"board_id" validate:"required"`
}

type MoveTaskRequest struct {
	ID           string `json:"id" validate:"required"`
	ColumnID     string `json:"column_id" validate:"required"`
	AfterTaskID  string `json:"after_task_id"`
	BeforeTaskID string `json:"before_task_id"`
	BoardID      string `json:"board_id" validate:"required"`
}

type UpdateTaskRequest struct {
	ID          string `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	AssigneeID  string `json:"assignee_id" validate:"required"`
	BoardID     string `json:"board_id" validate:"required"`
}

type DeleteTaskRequest struct {
	ID      string `json:"id" validate:"required"`
	BoardID string `json:"board_id" validate:"required"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	ColumnID    string    `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rank        string    `json:"rank"`
	AssigneeID  string    `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeleteTaskResponse struct{}
