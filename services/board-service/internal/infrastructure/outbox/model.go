package outbox

import (
	"time"

	"github.com/google/uuid"
)

// OutboxEvent представляет событие в таблице outbox для реализации паттерна Outbox.
type OutboxEvent struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	OccurredAt    time.Time
}
