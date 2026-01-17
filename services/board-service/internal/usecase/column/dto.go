package column

import "time"

type CreateColumnRequest struct {
	BoardID string `json:"board_id" validate:"required"`
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Rank    string `json:"rank" validate:"required"`
}

type MoveColumnRequest struct {
	ID             string `json:"id" validate:"required"`
	AfterColumnID  string `json:"after_column_id"`
	BeforeColumnID string `json:"before_column_id"`
}

type UpdateColumnRequest struct {
	ID    string `json:"id" validate:"required"`
	Title string `json:"title" validate:"required,min=1,max=255"`
}

type DeleteColumnRequest struct {
	ID string `json:"id" validate:"required"`
}

type ColumnResponse struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	Title     string    `json:"title"`
	Rank      string    `json:"rank"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteColumnResponse struct{}
