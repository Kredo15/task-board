package postgres

import (
	"time"
)

// boardModel — как доска выглядит в таблице Postgres
type boardModel struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	OwnerID     string    `db:"owner_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"update_at"`
}

// taskdModel — как задача выглядит в таблице Postgres
type taskdModel struct {
	ID          string    `db:"id"`
	ColumnID    string    `db:"column_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Rank        string    `db:"rank"`
	AssigneeID  string    `db:"assignee_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"update_at"`
}
