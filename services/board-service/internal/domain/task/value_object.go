package task

import (
	"unicode/utf8"
)

type TaskID string

func NewTaskID(id string) (TaskID, error) {
	if id == "" {
		return "", ErrInvalidTaskID
	}
	return TaskID(id), nil
}

type IDGenerator interface {
	Generate() string
}

type Title string

func NewTitle(v string) (Title, error) {
	if v == "" {
		return "", ErrInvalidTaskTitleEmpty
	}
	if utf8.RuneCountInString(v) > 100 {
		return "", ErrInvalidTaskTitleLong
	}
	return Title(v), nil
}

type Description string

func NewDescription(v string) (Description, error) {
	if utf8.RuneCountInString(v) > 1000 {
		return "", ErrInvaliTaskDescriptionLong
	}
	return Description(v), nil
}

type Rank string

func NewRank(r string) (Rank, error) {
	if r == "" {
		return "", ErrInvalidRank
	}
	return Rank(r), nil
}

type AssigneeID string

func NewAssigneeID(id string) (AssigneeID, error) {
	if id == "" {
		return "", ErrInvalidAssigneeID
	}
	return AssigneeID(id), nil
}
