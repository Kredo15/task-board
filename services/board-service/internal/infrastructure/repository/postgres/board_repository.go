package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (r *BoardRepository) Create(ctx context.Context, b *board.Board, event *board.DomainEvent) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Commit отменит этот вызов, если все пройдет успешно

	query := `
        INSERT INTO boards (id, title, description, owner_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err = tx.Exec(ctx, query,
		b.ID(),
		b.Title(),
		pgtype.Text{String: b.Description(), Valid: b.Description() != ""},
		b.OwnerID(),
		b.CreatedAt(),
		b.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert board: %w", err)
	}
	// Сохраняем событие в outbox, если оно есть
	if event != nil {
		const outboxQuery = `
            INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload, occurred_at)
            VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, outboxQuery,
			event.ID,
			"board",
			b.ID(),
			event.Type,
			event.Payload,
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

func (r *BoardRepository) Update(ctx context.Context, b *board.Board, event *board.DomainEvent) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Commit отменит этот вызов, если все пройдет успешно

	query := `
		UPDATE boards
		SET title = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := tx.Exec(ctx, query,
		b.Title(),
		pgtype.Text{String: b.Description(), Valid: b.Description() != ""},
		b.UpdatedAt(),
		b.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update board: %w", err)
	}

	if result.RowsAffected() == 0 {
		return board.ErrBoardNotFound
	}

	// Сохраняем событие в outbox, если оно есть
	if event != nil {
		const outboxQuery = `
            INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload, occurred_at)
            VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, outboxQuery,
			event.ID,
			"board",
			b.ID(),
			event.Type,
			event.Payload,
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

func (r *BoardRepository) Delete(ctx context.Context, id board.BoardID, event *board.DomainEvent) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Commit отменит этот вызов, если все пройдет успешно

	query := `
		DELETE FROM boards
		WHERE id = $1
	`
	result, err := tx.Exec(ctx, query, string(id))

	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	if result.RowsAffected() == 0 {
		return board.ErrBoardNotFound
	}

	// Сохраняем событие в outbox, если оно есть
	if event != nil {
		const outboxQuery = `
            INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload, occurred_at)
            VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, outboxQuery,
			event.ID,
			"board",
			string(id),
			event.Type,
			event.Payload,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, board.ErrBoardNotFound
		}
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	b := board.RestoreBoard(
		bModel.ID,
		bModel.Title,
		bModel.Description.String,
		bModel.OwnerID,
		bModel.CreatedAt,
		bModel.UpdatedAt,
		nil,
	)
	return b, nil
}

func (r *BoardRepository) GetFullBoard(ctx context.Context, id board.BoardID) (*board.Board, error) {
	// Запрос для получения доски с колонками и задачами
	query := `
        SELECT 
            b.id, b.title, b.description, b.owner_id, b.created_at, b.updated_at,
            c.id, c.title, c.rank, c.created_at, c.updated_at,
            t.id, t.column_id, t.title, t.content, t.rank, t.assignee_id, t.created_at, t.updated_at
        FROM boards b
        LEFT JOIN columns c ON b.id = c.board_id
        LEFT JOIN tasks t ON c.id = t.column_id
        WHERE b.id = $1
        ORDER BY c.rank, t.rank
    `

	rows, err := r.db.Query(ctx, query, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var b *board.Board
	columns := make(map[string]*board.Column)
	tasksMap := make(map[string][]*board.Task)

	for rows.Next() {
		var (
			// Временные переменные для сканирования
			bID, bTitle, bOwnerID  string
			bDesc                  pgtype.Text
			bCreatedAt, bUpdatedAt time.Time

			cID, cTitle, cRank     sql.NullString
			cCreatedAt, cUpdatedAt sql.NullTime

			tID, tColID, tTitle, tDesc, tRank, tAssigneeID sql.NullString
			tCreatedAt, tUpdatedAt                         sql.NullTime
		)

		err := rows.Scan(
			&bID, &bTitle, &bDesc, &bOwnerID, &bCreatedAt, &bUpdatedAt,
			&cID, &cTitle, &cRank, &cCreatedAt, &cUpdatedAt,
			&tID, &tColID, &tTitle, &tDesc, &tRank, &tAssigneeID, &tCreatedAt, &tUpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Инициализируем доску, если она еще не создана
		if b == nil {
			b = board.RestoreBoard(bID, bTitle, bDesc.String, bOwnerID, bCreatedAt, bUpdatedAt, nil)
		}
		// Обрабатываем колонку, если она существует
		if cID.Valid {
			if _, ok := columns[cID.String]; !ok {
				columns[cID.String], err = board.RestoreColumn(
					cID.String, bID, cTitle.String, cRank.String,
					cCreatedAt.Time, cUpdatedAt.Time, nil,
				)
				if err != nil {
					return nil, err
				}
			}
		}

		// Обрабатываем задачу, если она существует
		if tID.Valid {
			task, err := board.RestoreTask(
				tID.String, bID, tColID.String, tTitle.String,
				tDesc.String, tRank.String, tAssigneeID.String, tCreatedAt.Time, tUpdatedAt.Time,
			)
			if err != nil {
				return nil, err
			}
			tasksMap[tColID.String] = append(tasksMap[tColID.String], task)
		}
	}

	finalColumns := make([]*board.Column, 0, len(columns))
	for colID, col := range columns {
		colTasks := tasksMap[colID]
		col.SetTasks(colTasks)
		finalColumns = append(finalColumns, col)
	}

	b.SetColumns(finalColumns)

	return b, nil
}
