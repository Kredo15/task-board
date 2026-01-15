package board

import "time"

// CreateBoardRequest - запрос на создание доски
type CreateBoardRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	OwnerID     string `json:"owner_id" validate:"required"`
}

type GetBoardRequest struct {
	ID string `json:"id" validate:"required"`
}

type GetBoardsRequest struct {
	OwnerID string `json:"owner_id" validate:"required"`
}

type BoardResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetBoardsResponse struct {
	Boards []BoardResponse `json:"boards"`
}
