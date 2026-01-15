package task

import "errors"

var (
	ErrInvalidTaskID             = errors.New("task ID is required")
	ErrInvalidTaskTitleEmpty     = errors.New("task title cannot be empty")
	ErrInvalidTaskTitleLong      = errors.New("task title is too long")
	ErrInvaliTaskDescriptionLong = errors.New("task description is too long")
	ErrInvalidRank               = errors.New("invalid rank")
	ErrInvalidAssigneeID         = errors.New("invalid assignee id")
	ErrTaskNotFound              = errors.New("task not found")
)
