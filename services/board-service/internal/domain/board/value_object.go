package board

import (
	"strings"
	"unicode/utf8"
)

type BoardID string

func NewBoardID(id string) (BoardID, error) {
	if id == "" {
		return "", ErrInvalidBoardID
	}
	return BoardID(id), nil
}

type IDGenerator interface {
	Generate() string
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
