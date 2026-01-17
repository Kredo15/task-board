package column

import "errors"

var (
	ErrInvalidColumnID         = errors.New("column ID is required")
	ErrInvalidColumnTitleEmpty = errors.New("column title cannot be empty")
	ErrInvalidColumnTitleLong  = errors.New("column title is too long")
	ErrInvalidRank             = errors.New("invalid rank")
	ErrColumnNotFound          = errors.New("column not found")
)
