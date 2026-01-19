package outbox

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Relay struct {
	db     *pgxpool.Pool
	broker MessageBroker
	limit  int
}

func NewRelay(db *pgxpool.Pool, broker MessageBroker, limit int) *Relay {
	return &Relay{db: db, broker: broker, limit: limit}
}

// Start запускает бесконечный цикл обработки
func (r *Relay) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := r.processEvents(ctx); err != nil {
					log.Printf("Relay error: %v", err)
				}
			}
		}
	}()
}

func (r *Relay) processEvents(ctx context.Context) error {
	// 1. Начинаем транзакцию для обработки пачки (batch)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 2. Выбираем необработанные события и блокируем их (FOR UPDATE SKIP LOCKED)
	// Это позволяет запускать несколько реплик Relay одновременно без дублирования
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, occurred_at
		FROM outbox
		WHERE processed_at IS NULL
		ORDER BY occurred_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`
	rows, err := tx.Query(ctx, query, r.limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		if err := rows.Scan(&e.ID, &e.AggregateType, &e.AggregateID, &e.EventType, &e.Payload, &e.OccurredAt); err != nil {
			return err
		}
		events = append(events, e)
	}

	if len(events) == 0 {
		return nil
	}

	// 3. Отправляем события в брокер
	for _, e := range events {
		if err := r.broker.Publish(ctx, e); err != nil {
			// Если брокер недоступен, мы прерываем транзакцию.
			// События останутся в базе и будут обработаны в следующей итерации.
			return err
		}

		// 4. Помечаем как обработанное
		_, err = tx.Exec(ctx, "UPDATE outbox SET processed_at = $1 WHERE id = $2", time.Now(), e.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
