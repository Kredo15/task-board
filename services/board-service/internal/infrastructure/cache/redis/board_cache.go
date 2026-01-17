package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Kredo15/task-board/services/board-service/internal/usecase/board"
)

type BoardCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewBoardCache(client *redis.Client, ttl time.Duration) *BoardCache {
	return &BoardCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *BoardCache) Get(ctx context.Context, key string) (*board.BoardResponse, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var dto board.BoardResponse
	err = json.Unmarshal([]byte(val), &dto)
	return &dto, nil
}

func (c *BoardCache) Set(ctx context.Context, dto *board.BoardResponse) error {
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	err = c.client.Set(ctx, dto.ID, data, c.ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *BoardCache) Invalidate(ctx context.Context, id string) error {
	return c.client.Del(ctx, id).Err()
}
