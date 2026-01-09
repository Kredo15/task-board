package task

import (
	"unicode/utf8"
)

type TaskID string

func NewTaskID(id string) TaskID {
	return TaskID(id)
}

type IDGenerator interface {
	Generate() string
}

type Title string

func NewTitle(v string) (Title, error) {
	if v == "" {
		return "", ErrInvalidColumnTitleEmpty
	}
	if utf8.RuneCountInString(v) > 100 {
		return "", ErrInvalidColumnTitleLong
	}
	return Title(v), nil
}

type Description string

func NewDescription(v string) (Description, error) {
	if utf8.RuneCountInString(v) > 1000 {
		return "", ErrInvaliColumnDescriptionLong
	}
	return Description(v), nil
}

type Position int

func NewPosition(pos int) (Position, error) {
	if pos < 0 {
		return 0, ErrInvalidPosition
	}
	return Position(pos), nil
}

type AssigneeID string

func NewAssigneeID(id string) (AssigneeID, error) {
	if id == "" {
		return "", ErrInvalidAssigneeID
	}
	return AssigneeID(id), nil
}
