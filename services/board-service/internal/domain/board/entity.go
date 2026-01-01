package board

import (
	"time"
)

type Board struct {
	id          BoardID
	title       Title
	description Description
	ownerID     OwnerID
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBoard(gen IDGenerator, titleRaw, descRaw, ownerRaw string) (*Board, error) {
	// Валидируем Title
	title, err := NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем Description
	desc, err := NewDescription(descRaw)
	if err != nil {
		return nil, err
	}
	// Валидируем OwnerID
	owner_id, err := NewOwnerID(ownerRaw)
	if err != nil {
		return nil, err
	}

	board := &Board{
		id:          BoardID(gen.Generate()),
		title:       title,
		description: desc,
		ownerID:     owner_id,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}

	return board, nil
}

func RestoreBoard(id, title, desc, ownerID string, createdAt, updatedAt time.Time) *Board {
	return &Board{
		id:          BoardID(id),
		title:       Title(title),
		description: Description(desc),
		ownerID:     OwnerID(ownerID),
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (b *Board) ID() string { return string(b.id) }

func (b *Board) Title() string { return string(b.title) }

func (b *Board) Description() string { return string(b.description) }

func (b *Board) OwnerID() string { return string(b.ownerID) }

func (b *Board) CreatedAt() time.Time { return b.createdAt }

func (b *Board) UpdatedAt() time.Time { return b.updatedAt }

func (b *Board) UpdateTitle(newTitleRaw string) error {
	title, err := NewTitle(newTitleRaw)
	if err != nil {
		return err
	}
	b.title = title
	b.updatedAt = time.Now()
	return nil
}

func (b *Board) UpdateDescription(newDescRaw string) error {
	desc, err := NewDescription(newDescRaw)
	if err != nil {
		return err
	}
	b.description = desc
	b.updatedAt = time.Now()
	return nil
}

func (b *Board) Equals(other *Board) bool {
	if other == nil {
		return false
	}
	return b.id == other.id
}
