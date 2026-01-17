package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *BoardRepository) Create(ctx context.Context, b *board.Board) error {

	query := `
        INSERT INTO boards (id, title, description, owner_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
		RETURNIG ID
    `

	err := r.db.QueryRow(ctx, query,
		b.ID(),
		b.Title(),
		pgtype.Text{String: b.Description(), Valid: true},
		b.OwnerID(),
		b.CreatedAt(),
		b.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to create board: %s", err)
	}
	return nil
}

func (r *BoardRepository) GetByID(ctx context.Context, id board.BoardID) (*board.Board, error) {
	var bModel boardModel

	query := `
		SELECT id, title, description, owner_id, created_at, updated_at
		FROM boards
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, string(id)).Scan(
		&bModel.ID,
		&bModel.Title,
		&bModel.Description,
		&bModel.OwnerID,
		&bModel.CreatedAt,
		&bModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, board.ErrBoardNotFound
		}
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	b := board.RestoreBoard(
		bModel.ID,
		bModel.Title,
		bModel.Description,
		bModel.OwnerID,
		bModel.CreatedAt,
		bModel.UpdatedAt,
	)
	return b, nil
}

func (r *BoardRepository) GetBoards(ctx context.Context, owner_id board.OwnerID) ([]*board.Board, error) {
	query := `
		SELECT id, title, description, owner_id, created_at, updated_at
		FROM boards
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, owner_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]*board.Board, 0)

	for rows.Next() {
		var bModel boardModel

		err := rows.Scan(
			&bModel.ID,
			&bModel.Title,
			&bModel.Description,
			&bModel.OwnerID,
			&bModel.CreatedAt,
			&bModel.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		b := board.RestoreBoard(
			bModel.ID,
			bModel.Title,
			bModel.Description,
			bModel.OwnerID,
			bModel.CreatedAt,
			bModel.UpdatedAt,
		)

		boards = append(boards, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return boards, nil
}

func (r *BoardRepository) GetFullBoard(ctx context.Context, id board.BoardID) (*board.Board, error) {
	// Запрос для получения доски с колонками и задачами
	/*query := `
	        SELECT
				b.id, b.title, b.description, b.owner_id,
				c.id, c.title, c.rank,
				t.id, t.title, t.description, t.rank, t.assignee_id
	        FROM boards b
	        LEFT JOIN columns c ON b.id = c.board_id
	        LEFT JOIN tasks t ON c.id = t.column_id
	        WHERE b.id = $1
	        ORDER BY c.rank, t.rank
	`

		rows, err := r.db.Query(ctx, query, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var b *board.Entity
		// Используем map, чтобы не дублировать колонки
		// key: column_id, value: индекс в слайсе b.Columns
		colMap := make(map[string]*column.Entity)

		for rows.Next() {
			var row boardRow
			err := rows.Scan(
				&row.BoardID, &row.BoardTitle,
				&row.ColID, &row.ColTitle, &row.ColPos,
				&row.TaskID, &row.TaskTitle, &row.TaskPos,
			)
			if err != nil {
				return nil, err
			}

			// 1. Инициализируем доску (только на первой итерации)
			if b == nil {
				b = &board.Entity{
					ID:      row.BoardID,
					Title:   row.BoardTitle,
					Columns: []*column.Entity{},
				}
			}

			// 2. Обрабатываем колонку (если она есть)
			if row.ColID != nil {
				col, exists := colMap[*row.ColID]
				if !exists {
					col = &column.Entity{
						ID:       *row.ColID,
						Title:    *row.ColTitle,
						Position: *row.ColPos,
						Tasks:    []*task.Entity{},
					}
					b.Columns = append(b.Columns, col)
					colMap[*row.ColID] = col
				}

				// 3. Обрабатываем задачу (если она есть в этой колонке)
				if row.TaskID != nil {
					t := &task.Entity{
						ID:       *row.TaskID,
						Title:    *row.TaskTitle,
						Position: *row.TaskPos,
					}
					col.Tasks = append(col.Tasks, t)
				}
			}
		}

		if b == nil {
			return nil, board.ErrBoardNotFound
		}

		return b, nil */
	return nil, nil
}

func (r *BoardRepository) Update(ctx context.Context, b *board.Board) error {
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
