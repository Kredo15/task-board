package task

import "errors"

var (
	ErrInvalidColumnTitleEmpty     = errors.New("column title cannot be empty")
	ErrInvalidColumnTitleLong      = errors.New("column title is too long")
	ErrInvaliColumnDescriptionLong = errors.New("column description is too long")
	ErrInvalidPosition             = errors.New("invalid positioin")
	ErrInvalidAssigneeID           = errors.New("invalid assignee id")
)
