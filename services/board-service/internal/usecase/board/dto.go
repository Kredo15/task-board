package board

import "time"

// CreateBoardRequest - запрос на создание доски
type CreateBoardRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	OwnerID     string `json:"owner_id" validate:"required"`
}

type CreateBoardResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}
