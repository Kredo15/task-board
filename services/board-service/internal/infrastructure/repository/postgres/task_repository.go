package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/column"
	"github.com/Kredo15/task-board/services/board-service/internal/domain/task"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {

	query := `
        INSERT INTO tasks (
		id, column_id, title, description, position, assignee_id, created_at, updated_at
		)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNIG ID
    `

	err := r.db.QueryRow(ctx, query,
		t.ID(),
		t.ColumnID(),
		t.Title(),
		pgtype.Text{String: t.Description(), Valid: true},
		t.Rank(),
		t.AssigneeID(),
		t.CreatedAt(),
		t.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to create task: %s", err)
	}
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id task.TaskID) (*task.Task, error) {
	var t taskdModel
	query := `
		SELECT id, column_id, title, description, rank, assignee_id, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, string(id)).Scan(
		&t.ID,
		&t.ColumnID,
		&t.Title,
		&t.Description,
		&t.Rank,
		&t.AssigneeID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, task.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	task := task.RestoreTask(
		t.ID,
		t.ColumnID,
		t.Title,
		t.Description,
		t.Rank,
		t.AssigneeID,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return task, nil

}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	query := `
        UPDATE tasks
        SET title = $1, description = $2, updated_at = $3
        WHERE id = $4
    `

	result, err := r.db.Exec(ctx, query,
		t.Title(),
		pgtype.Text{String: t.Description(), Valid: true},
		t.UpdatedAt(),
		t.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update board: %w", err)
	}

	if result.RowsAffected() == 0 {
		return task.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) GetRanksByColumn(ctx context.Context, colID column.ColumnID) ([]string, error) {
	var ranks []string
	query := `SELECT rank FROM tasks WHERE column_id = $1 ORDER BY rank ASC`
	rows, err := r.db.Query(ctx, query, colID)
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

func (r *TaskRepository) GetIndexByRank(ctx context.Context, colID column.ColumnID, rank task.Rank) (int32, error) {
	var count int32
	// Считаем, сколько задач имеют ранг меньше текущего
	query := `SELECT COUNT(*) FROM tasks WHERE column_id = $1 AND rank < $2`

	err := r.db.QueryRow(ctx, query, colID, rank).Scan(&count)

	return count, err
}

func (r *TaskRepository) Delete(ctx context.Context, id task.TaskID) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, string(id))

	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	if result.RowsAffected() == 0 {
		return task.ErrTaskNotFound
	}
	return nil
}
