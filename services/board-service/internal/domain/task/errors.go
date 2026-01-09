package task

import "errors"

var (
	ErrInvalidTaskID               = errors.New("task ID is required")
	ErrInvalidColumnTitleEmpty     = errors.New("task title cannot be empty")
	ErrInvalidColumnTitleLong      = errors.New("task title is too long")
	ErrInvaliColumnDescriptionLong = errors.New("task description is too long")
	ErrInvalidPosition             = errors.New("invalid positioin")
	ErrInvalidAssigneeID           = errors.New("invalid assignee id")
)
