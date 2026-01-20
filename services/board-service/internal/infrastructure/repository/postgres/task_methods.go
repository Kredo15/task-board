package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

func (r *BoardRepository) SaveTask(
	ctx context.Context,
	boardID board.BoardID,
	t *board.Task,
	event *board.DomainEvent,
) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Commit отменит этот вызов, если все пройдет успешно

	query := `
		INSERT INTO tasks (id, board_id, column_id, title, content, rank, assignee_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			column_id = EXCLUDED.column_id,
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			rank = EXCLUDED.rank,
			assignee_id = EXCLUDED.assignee_id,
			updated_at = EXCLUDED.updated_at
	`
	_, err = tx.Exec(ctx, query,
		t.ID(),
		string(boardID),
		t.ColumnID(),
		t.Title(),
		pgtype.Text{String: t.Description(), Valid: t.Description() != ""},
		t.Rank(),
		t.AssigneeID(),
		t.CreatedAt(),
		t.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// Сохраняем событие в outbox, если оно есть
	if event != nil {
		// Сериализуем полезную нагрузку события в JSON
		payloadBytes, err := json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}

		const outboxQuery = `
			INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload, occurred_at)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.Exec(ctx, outboxQuery,
			event.ID,
			"board",
			string(boardID),
			event.Type,
			payloadBytes,
			event.OccurredAt,
		)
		if err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *BoardRepository) GetTaskByID(
	ctx context.Context,
	boardID board.BoardID,
	taskID board.TaskID,
) (*board.Task, error) {
	query := `
		SELECT id, board_id, column_id, title, content, rank, assignee_id, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND board_id = $2
	`
	row := r.db.QueryRow(ctx, query, string(taskID), string(boardID))
	var tModel taskModel
	err := row.Scan(
		&tModel.ID,
		&tModel.BoardID,
		&tModel.ColumnID,
		&tModel.Title,
		&tModel.Description,
		&tModel.Rank,
		&tModel.AssigneeID,
		&tModel.CreatedAt,
		&tModel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	desc := ""
	if tModel.Description.Valid {
		desc = tModel.Description.String
	}
	t, err := board.RestoreTask(
		tModel.ID,
		tModel.BoardID,
		tModel.ColumnID,
		tModel.Title,
		desc,
		tModel.Rank,
		tModel.AssigneeID,
		tModel.CreatedAt,
		tModel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to restore task: %w", err)
	}
	return t, nil
}

func (r *BoardRepository) DeleteTask(
	ctx context.Context,
	boardID board.BoardID,
	taskID board.TaskID,
	event *board.DomainEvent,
) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Commit отменит этот вызов, если все пройдет успешно

	query := `
		DELETE FROM tasks
		WHERE id = $1 AND board_id = $2
	`
	result, err := tx.Exec(ctx, query, string(taskID), string(boardID))
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	if result.RowsAffected() == 0 {
		return board.ErrTaskNotFound
	}
	// Сохраняем событие в outbox, если оно есть
	if event != nil {
		// Сериализуем полезную нагрузку события в JSON
		payloadBytes, err := json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}

		const outboxQuery = `
			INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload, occurred_at)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.Exec(ctx, outboxQuery,
			event.ID,
			"board",
			string(boardID),
			event.Type,
			payloadBytes,
			event.OccurredAt,
		)
		if err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *BoardRepository) GetTaskRanks(ctx context.Context, bID board.BoardID, colID board.ColumnID) ([]board.TaskRank, error) {
	query := `
		SELECT id, rank
		FROM tasks
		WHERE column_id = $1 AND board_id = $2
		ORDER BY rank ASC
	`
	rows, err := r.db.Query(ctx, query, string(colID), string(bID))
	if err != nil {
		return nil, fmt.Errorf("failed to get task ranks: %w", err)
	}
	defer rows.Close()

	var ranks []board.TaskRank
	for rows.Next() {
		var rank board.TaskRank
		if err := rows.Scan(&rank.ID, &rank.Rank); err != nil {
			return nil, fmt.Errorf("failed to scan rank: %w", err)
		}
		ranks = append(ranks, rank)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return ranks, nil
}
