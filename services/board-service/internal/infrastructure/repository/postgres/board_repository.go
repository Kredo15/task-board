package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

type BoardRepository struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(ctx context.Context, board *board.Board) error {

	query := `
        INSERT INTO boards (title, description, owner_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
		RETURNIG ID
    `

	err := r.db.QueryRow(ctx, query,
		board.ID(),
		board.Title(),
		pgtype.Text{String: board.Description(), Valid: true},
		board.OwnerID(),
		board.CreatedAt(),
		board.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to create board: %s", err)
	}
	return nil
}

func (r *BoardRepository) GetBoard(ctx context.Context, id board.BoardID) (*board.Board, error) {
	var (
		rawID       string
		title       string
		description string
		ownerID     string
		createdAt   time.Time
		updatedAt   time.Time
	)

	query := `
		SELECT id, title, description, owner_id, created_at, updated_at
		FROM boards
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, string(id)).Scan(
		&rawID,
		&title,
		&description,
		&ownerID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, board.ErrBoardNotFound
		}
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	b := board.RestoreBoard(
		rawID,
		title,
		description,
		ownerID,
		createdAt,
		updatedAt,
	)
	return b, nil
}

func (r *BoardRepository) GetBoards(ctx context.Context) ([]*board.Board, error) {
	query := `
		SELECT id, title, description, owner_id, created_at, updated_at
		FROM boards
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]*board.Board, 0)

	for rows.Next() {
		var (
			rawID       string
			title       string
			description string
			ownerID     string
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(
			&rawID,
			&title,
			&description,
			&ownerID,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		b := board.RestoreBoard(
			rawID,
			title,
			description,
			ownerID,
			createdAt,
			updatedAt,
		)

		boards = append(boards, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return boards, nil
}

func (r *BoardRepository) Update(ctx context.Context, b board.Board) error {
	query := `
        UPDATE boards
        SET title = $1, description = $2, updated_at = $3
        WHERE id = $4
    `

	result, err := r.db.Exec(ctx, query,
		b.Title(),
		pgtype.Text{String: b.Description(), Valid: true},
		b.UpdatedAt(),
		b.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update board: %w", err)
	}

	if result.RowsAffected() == 0 {
		return board.ErrBoardNotFound
	}

	return nil
}

func (r *BoardRepository) Delete(ctx context.Context, id board.BoardID) error {
	query := `
		DELETE FROM boards
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, string(id))

	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	if result.RowsAffected() == 0 {
		return board.ErrBoardNotFound
	}
	return nil
}
