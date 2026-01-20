package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Kredo15/task-board/services/board-service/internal/domain/board"
)

func (r *BoardRepository) SaveColumn(ctx context.Context, boardID board.BoardID, c *board.Column, event *board.DomainEvent) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Commit отменит этот вызов, если все пройдет успешно

	// Сохраняем или обновляем колонку

	query := `
        INSERT INTO columns (id, board_id, title, rank, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE SET
            title = EXCLUDED.title,
            rank = EXCLUDED.rank,
            updated_at = EXCLUDED.updated_at
    `
	_, err = tx.Exec(ctx, query,
		c.ID(),
		string(boardID),
		c.Title(),
		c.Rank(),
		c.CreatedAt(),
		c.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to save column: %w", err)
	}

	// Сериализуем полезную нагрузку события в JSON
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	// Сохраняем событие в outbox, если оно есть
	if event != nil {
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

func (r *BoardRepository) GetColumnByID(ctx context.Context, boardID board.BoardID, columnID board.ColumnID) (*board.Column, error) {
	query := `
		SELECT id, board_id, title, rank, created_at, updated_at
		FROM columns
		WHERE id = $1 AND board_id = $2
	`
	row := r.db.QueryRow(ctx, query, string(columnID), string(boardID))
	var cModel columnModel
	err := row.Scan(&cModel.ID, &cModel.BoardID, &cModel.Title, &cModel.Rank, &cModel.CreatedAt, &cModel.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, board.ErrColumnNotFound
		}
		return nil, fmt.Errorf("failed to scan column: %w", err)
	}
	col, err := board.RestoreColumn(
		cModel.ID,
		cModel.BoardID,
		cModel.Title,
		cModel.Rank,
		cModel.CreatedAt,
		cModel.UpdatedAt,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create column domain model: %w", err)
	}
	return col, nil
}

func (r *BoardRepository) DeleteColumn(
	ctx context.Context,
	boardID board.BoardID,
	colID board.ColumnID,
	event *board.DomainEvent,
) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Commit отменит этот вызов, если все пройдет успешно

	query := `
		DELETE FROM columns
		WHERE id = $1 AND board_id = $2
	`
	result, err := tx.Exec(ctx, query, string(colID), string(boardID))
	if err != nil {
		return fmt.Errorf("failed to delete column: %w", err)
	}

	if result.RowsAffected() == 0 {
		return board.ErrColumnNotFound
	}

	// Сериализуем полезную нагрузку события в JSON
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	// Сохраняем событие в outbox, если оно есть
	if event != nil {
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

func (r *BoardRepository) GetColumnRanks(ctx context.Context, boardID board.BoardID) ([]board.ColumnRank, error) {
	query := `
		SELECT id, rank
		FROM columns
		WHERE board_id = $1
		ORDER BY rank ASC
	`
	rows, err := r.db.Query(ctx, query, string(boardID))
	if err != nil {
		return nil, fmt.Errorf("failed to get column ranks: %w", err)
	}
	defer rows.Close()

	var ranks []board.ColumnRank
	for rows.Next() {
		var rank board.ColumnRank
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
