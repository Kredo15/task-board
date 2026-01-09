package column

import (
	"unicode/utf8"
)

type ColumnID string

func NewColumnID(id string) (ColumnID, error) {
	if id == "" {
		return "", ErrInvalidColumnID
	}
	return ColumnID(id), nil
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

type Position int

func NewPosition(pos int) (Position, error) {
	if pos < 0 {
		return 0, ErrInvalidPosition
	}
	return Position(pos), nil
}
