package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/Kredo15/task-board/services/board-service/internal/infrastructure/outbox"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: writer}
}

func (kp *Producer) Publish(ctx context.Context, event outbox.OutboxEvent) error {
	// Формируем сообщение Kafka
	msg := kafka.Message{
		Key:   []byte(event.AggregateID), // Группируем сообщения по ID агрегата
		Value: event.Payload,
		Headers: []kafka.Header{ // Используем заголовки чтобы не парсить полезную нагрузку для определения типа события
			{Key: "event_type", Value: []byte(event.EventType)},
			{Key: "aggregate_type", Value: []byte(event.AggregateType)},
			{Key: "version", Value: []byte("v1")}, // Помогает потребителю выбрать схему декодирования
		},
	}
	err := kp.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message to Kafka: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
