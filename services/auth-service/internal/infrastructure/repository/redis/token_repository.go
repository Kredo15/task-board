package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Kredo15/task-board/services/auth-service/internal/domain"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) Save(ctx context.Context, token *domain.RefreshToken) error {
	ttl := time.Until(token.ExpiresAt)
	// Ключ для конкретного токена
	tokenKey := fmt.Sprintf("token:%s", token)
	// Ключ для списка всех токенов пользователя (индекс)
	userTokensKey := fmt.Sprintf("user_tokens:%s", token.UserID)

	// Используем Pipeline для атомарности и производительности
	pipe := r.client.Pipeline()

	// Сохраняем связку токен -> userID с TTL
	pipe.Set(ctx, tokenKey, token.UserID, ttl)

	// Добавляем токен в список токенов пользователя
	pipe.SAdd(ctx, userTokensKey, token)

	// Устанавливаем TTL для списка (чтобы мусор удалялся автоматически)
	pipe.Expire(ctx, userTokensKey, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save refresh token to redis: %w", err)
	}

	return nil
}

func (r *TokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	tokenKey := fmt.Sprintf("token:%s", token)

	// Выполняем запрос к Redis
	userID, err := r.client.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	return &domain.RefreshToken{
		Token:  token,
		UserID: userID,
	}, nil
}

func (r *TokenRepository) DeleteByToken(ctx context.Context, token string) error {
	tokenKey := fmt.Sprintf("token:%s", token)

	// Сначала узнаем, какому пользователю принадлежит токен
	// Это нужно, чтобы найти нужный Set (user_tokens:ID)
	userID, err := r.client.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil // Токена и так нет, считать успешным
		}
		return fmt.Errorf("failed to find token before deletion: %w", err)
	}

	userTokensKey := fmt.Sprintf("user_tokens:%s", userID)

	// Используем Pipeline для атомарного удаления
	pipe := r.client.Pipeline()

	// Удаляем основной ключ токена
	pipe.Del(ctx, tokenKey)

	// Удаляем токен из списка токенов этого пользователя
	pipe.SRem(ctx, userTokensKey, token)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute delete pipeline: %w", err)
	}

	return nil
}

func (r *TokenRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	userTokensKey := fmt.Sprintf("user_tokens:%s", userID)

	// Получаем все активные токены пользователя из индекса (Set)
	tokens, err := r.client.SMembers(ctx, userTokensKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil // Токенов нет, делать нечего
		}
		return fmt.Errorf("failed to fetch user tokens: %w", err)
	}

	if len(tokens) == 0 {
		return nil
	}

	// Используем Pipeline для массового удаления
	pipe := r.client.Pipeline()

	// Удаляем каждый индивидуальный ключ токена
	for _, token := range tokens {
		tokenKey := fmt.Sprintf("token:%s", token)
		pipe.Del(ctx, tokenKey)
	}

	// Удаляем сам индексный ключ пользователя
	pipe.Del(ctx, userTokensKey)

	// Выполняем все команды одним пакетом
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute mass delete pipeline: %w", err)
	}

	return nil
}
