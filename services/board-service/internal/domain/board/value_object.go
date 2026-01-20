package board

import (
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

type IDGenerator interface {
	Generate() string
}

type LexorankGen interface {
	Between(prevKey, nextKey string) (string, error)
}

type BoardID string

func ParseBoardID(id string) (BoardID, error) {
	if id == "" {
		return "", ErrInvalidBoardID
	}
	if _, err := uuid.Parse(id); err != nil {
		return "", ErrInvalidBoardID
	}

	return BoardID(id), nil
}

type Title string

func ParseTitle(v string) (Title, error) {
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

func ParseDescription(v string) (Description, error) {
	v = strings.TrimSpace(v)
	if utf8.RuneCountInString(v) > 1000 {
		return "", ErrInvalidBoardDescriptionLong
	}
	return Description(v), nil
}

type OwnerID string

func ParseOwnerID(id string) (OwnerID, error) {
	if id == "" {
		return "", ErrInvalidOwnerID
	}
	return OwnerID(id), nil
}

type ColumnID string

func ParseColumnID(id string) (ColumnID, error) {
	if id == "" {
		return "", ErrInvalidColumnID
	}
	if _, err := uuid.Parse(id); err != nil {
		return "", ErrInvalidColumnID
	}
	return ColumnID(id), nil
}

type Rank string

func ParseRank(r string) (Rank, error) {
	if r == "" {
		return "", ErrInvalidRank
	}
	return Rank(r), nil
}

type TaskID string

func ParseTaskID(id string) (TaskID, error) {
	if id == "" {
		return "", ErrInvalidTaskID
	}
	if _, err := uuid.Parse(id); err != nil {
		return "", ErrInvalidTaskID
	}
	return TaskID(id), nil
}

type AssigneeID string

func ParseAssigneeID(id string) (AssigneeID, error) {
	if id == "" {
		return "", ErrInvalidAssigneeID
	}
	if _, err := uuid.Parse(id); err != nil {
		return "", ErrInvalidAssigneeID
	}
	return AssigneeID(id), nil
}

type EventID string

func ParseEventID(id string) (EventID, error) {
	if id == "" {
		return "", ErrInvalidEventID
	}
	if _, err := uuid.Parse(id); err != nil {
		return "", ErrInvalidEventID
	}
	return EventID(id), nil
}
