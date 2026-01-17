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
	var t taskModel
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

	taskRestore, err := task.RestoreTask(
		t.ID,
		t.ColumnID,
		t.Title,
		t.Description,
		t.Rank,
		t.AssigneeID,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return taskRestore, nil

}
func (r *TaskRepository) GetByColumn(ctx context.Context, colID column.ColumnID) ([]*task.Task, error) {
	query := `
		SELECT id, column_id, title, description, rank, assignee_id, created_at, updated_at
		FROM tasks
		WHERE column_id = $1
		ORDER BY rank ASC
	`
	rows, err := r.db.Query(ctx, query, colID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*task.Task, 0)

	for rows.Next() {
		var tModel taskModel

		err := rows.Scan(
			&tModel.ID,
			&tModel.ColumnID,
			&tModel.Title,
			&tModel.Description,
			&tModel.Rank,
			&tModel.AssigneeID,
			&tModel.CreatedAt,
			&tModel.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		t, err := task.RestoreTask(
			tModel.ID,
			tModel.ColumnID,
			tModel.Title,
			tModel.Description,
			tModel.Rank,
			tModel.AssigneeID,
			tModel.CreatedAt,
			tModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetByColumns(ctx context.Context, colsID []column.ColumnID) ([]*task.Task, error) {
	if len(colsID) == 0 {
		return []*task.Task{}, nil
	}

	query := `
		SELECT id, column_id, title, description, rank, assignee_id, created_at, updated_at
		FROM tasks
		WHERE column_id = $1
		ORDER BY rank ASC
	`
	rows, err := r.db.Query(ctx, query, colsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*task.Task, 0)

	for rows.Next() {
		var tModel taskModel

		err := rows.Scan(
			&tModel.ID,
			&tModel.ColumnID,
			&tModel.Title,
			&tModel.Description,
			&tModel.Rank,
			&tModel.AssigneeID,
			&tModel.CreatedAt,
			&tModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		t, err := task.RestoreTask(
			tModel.ID,
			tModel.ColumnID,
			tModel.Title,
			tModel.Description,
			tModel.Rank,
			tModel.AssigneeID,
			tModel.CreatedAt,
			tModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
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

func (r *TaskRepository) GetRanksByID(ctx context.Context, afterTaskID, beforeTaskID string) ([]string, error) {
	ranks := make([]string, 0, 2)
	query := `SELECT rank FROM tasks WHERE task_id IN ($1, $2) ORDER BY rank ASC`
	rows, err := r.db.Query(ctx, query, afterTaskID, beforeTaskID)
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

func (r *TaskRepository) GetMaxRanksByColumn(ctx context.Context, colID column.ColumnID) (string, error) {
	var rank string
	query := `SELECT rank FROM tasks WHERE column_id = $1 ORDER BY rank DESC LIMIT 1`
	err := r.db.QueryRow(ctx, query, string(colID)).Scan(&rank)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return rank, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id task.TaskID) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, string(id))

	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return task.ErrTaskNotFound
	}
	return nil
}
