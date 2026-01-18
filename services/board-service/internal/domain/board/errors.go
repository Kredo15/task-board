package board

import "errors"

var (
	ErrInvalidBoardID              = errors.New("board ID is required")
	ErrInvalidBoardTitleEmpty      = errors.New("board title cannot be empty")
	ErrInvalidBoardTitleLong       = errors.New("board title is too long")
	ErrInvalidBoardDescriptionLong = errors.New("board description is too long")
	ErrInvalidOwnerID              = errors.New("board ownerID is required")
	ErrBoardNotFound               = errors.New("board not found")

	ErrInvalidColumnID         = errors.New("column ID is required")
	ErrInvalidColumnTitleEmpty = errors.New("column title cannot be empty")
	ErrInvalidColumnTitleLong  = errors.New("column title is too long")
	ErrInvalidRank             = errors.New("invalid rank")
	ErrColumnNotFound          = errors.New("column not found")

	ErrInvalidTaskID             = errors.New("task ID is required")
	ErrInvalidTaskTitleEmpty     = errors.New("task title cannot be empty")
	ErrInvalidTaskTitleLong      = errors.New("task title is too long")
	ErrInvaliTaskDescriptionLong = errors.New("task description is too long")
	ErrInvalidAssigneeID         = errors.New("invalid assignee id")
	ErrTaskNotFound              = errors.New("task not found")
	ErrLimitTaskReached          = errors.New("task limit reached for the column")
)
