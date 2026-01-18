package board

import (
	"strings"
	"unicode/utf8"
)

type IDGenerator interface {
	Generate() string
}

type LexorankGen interface {
	Between(prevKey, nextKey string) (string, error)
}

type BoardID string

func NewBoardID(id string) (BoardID, error) {
	if id == "" {
		return "", ErrInvalidBoardID
	}
	return BoardID(id), nil
}

type Title string

func NewTitle(v string) (Title, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", ErrInvalidBoardTitleEmpty
	}
	if utf8.RuneCountInString(v) > 100 {
		return "", ErrInvalidBoardTitleLong
	}
	return Title(v), nil
}

type Description string

func NewDescription(v string) (Description, error) {
	v = strings.TrimSpace(v)
	if utf8.RuneCountInString(v) > 1000 {
		return "", ErrInvalidBoardDescriptionLong
	}
	return Description(v), nil
}

type OwnerID string

func NewOwnerID(id string) (OwnerID, error) {
	if id == "" {
		return "", ErrInvalidOwnerID
	}
	return OwnerID(id), nil
}

type ColumnID string

func NewColumnID(id string) (ColumnID, error) {
	if id == "" {
		return "", ErrInvalidColumnID
	}
	return ColumnID(id), nil
}

type Rank string

func NewRank(r string) (Rank, error) {
	if r == "" {
		return "", ErrInvalidRank
	}
	return Rank(r), nil
}

type TaskID string

func NewTaskID(id string) (TaskID, error) {
	if id == "" {
		return "", ErrInvalidTaskID
	}
	return TaskID(id), nil
}

type AssigneeID string

func NewAssigneeID(id string) (AssigneeID, error) {
	if id == "" {
		return "", ErrInvalidAssigneeID
	}
	return AssigneeID(id), nil
}
