package board

import "errors"

var (
	ErrInvalidBoardID              = errors.New("board ID is required")
	ErrInvalidBoardTitleEmpty      = errors.New("board title cannot be empty")
	ErrInvalidBoardTitleLong       = errors.New("board title is too long")
	ErrInvalidBoardDescriptionLong = errors.New("board description is too long")
	ErrInvalidOwnerID              = errors.New("board ownerID is required")
	ErrBoardNotFound               = errors.New("board not found")
)
