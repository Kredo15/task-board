package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
)

type ColumnRepository struct {
	db *pgxpool.Pool
}

func NewColumnRepository(db *pgxpool.Pool) *ColumnRepository {
	return &ColumnRepository{db: db}
}

func (r *ColumnRepository) Create(ctx context.Context, c *column.Column) error {
	query := `
        INSERT INTO columns (
		id, board_id, title, created_at, updated_at
		)
        VALUES ($1, $2, $3, $4, $5)
		RETURNIG ID
    `

	err := r.db.QueryRow(ctx, query,
		c.ID(),
		c.BoardID(),
		c.Title(),
		c.CreatedAt(),
		c.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to create column: %s", err)
	}
	return nil
}

func (r *ColumnRepository) GetByID(ctx context.Context, id column.ColumnID) (*column.Column, error) {
	var t columnModel
	query := `
		SELECT id, board_id, title, rank, created_at, updated_at
		FROM columns
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, string(id)).Scan(
		&t.ID,
		&t.BoardID,
		&t.Title,
		&t.Rank,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, column.ErrColumnNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	columnRestore, err := column.RestoreColumn(
		t.ID,
		t.BoardID,
		t.Title,
		t.Rank,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return columnRestore, nil
}

func (r *ColumnRepository) Update(ctx context.Context, c *column.Column) error {
	query := `
        UPDATE tasks
        SET 
			board_id = $1, 
			title = $2, 
			rank = $3, 
			updated_at = $4
        WHERE id = $5
    `

	result, err := r.db.Exec(ctx, query,
		c.BoardID(),
		c.Title(),
		c.Rank(),
		c.UpdatedAt(),
		c.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update column: %w", err)
	}

	if result.RowsAffected() == 0 {
		return column.ErrColumnNotFound
	}

	return nil
}

func (r *ColumnRepository) GetByBoard(ctx context.Context, boardID board.BoardID) ([]*column.Column, error) {
	query := `
		SELECT id, board_id, title, rank, created_at, updated_at
		FROM columns
		WHERE board_id = $1
		ORDER BY rank ASC
	`
	rows, err := r.db.Query(ctx, query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]*column.Column, 0)

	for rows.Next() {
		var cModel columnModel

		err := rows.Scan(
			&cModel.ID,
			&cModel.BoardID,
			&cModel.Title,
			&cModel.Rank,
			&cModel.CreatedAt,
			&cModel.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		c, err := column.RestoreColumn(
			cModel.ID,
			cModel.BoardID,
			cModel.Title,
			cModel.Rank,
			cModel.CreatedAt,
			cModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		columns = append(columns, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

func (r *ColumnRepository) GetByBoards(ctx context.Context, boardsID []board.BoardID) ([]*column.Column, error) {
	if len(boardsID) == 0 {
		return []*column.Column{}, nil
	}

	query := `
		SELECT id, board_id, title, rank, created_at, updated_at
		FROM columns
		WHERE board_id = ANY($1)
		ORDER BY board_id ASC, rank ASC
	`

	rows, err := r.db.Query(ctx, query, boardsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]*column.Column, 0)

	for rows.Next() {
		var cModel columnModel

		err := rows.Scan(
			&cModel.ID,
			&cModel.BoardID,
			&cModel.Title,
			&cModel.Rank,
			&cModel.CreatedAt,
			&cModel.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		c, err := column.RestoreColumn(
			cModel.ID,
			cModel.BoardID,
			cModel.Title,
			cModel.Rank,
			cModel.CreatedAt,
			cModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		columns = append(columns, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

func (r *ColumnRepository) GetRanksByID(ctx context.Context, afterColumnID, beforeColumnID string) ([]string, error) {
	ranks := make([]string, 0, 2)
	query := `SELECT rank FROM tasks WHERE task_id IN ($1, $2) ORDER BY rank ASC`
	rows, err := r.db.Query(ctx, query, afterColumnID, beforeColumnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rank string
		err := rows.Scan(&rank)
		if err != nil {
			return nil, err
		}
		ranks = append(ranks, rank)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ranks, err
}

func (r *ColumnRepository) GetMaxRanksByBoard(ctx context.Context, boardID board.BoardID) (string, error) {
	var rank string
	query := `SELECT rank FROM columns WHERE board_id = $1 ORDER BY rank DESC LIMIT 1`
	err := r.db.QueryRow(ctx, query, string(boardID)).Scan(&rank)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return rank, nil
}

func (r *ColumnRepository) Delete(ctx context.Context, id column.ColumnID) error {
	query := `
		DELETE FROM columns
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, string(id))

	if err != nil {
		return fmt.Errorf("failed to delete column: %w", err)
	}

	if result.RowsAffected() == 0 {
		return column.ErrColumnNotFound
	}
	return nil
}
