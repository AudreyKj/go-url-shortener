package storage

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis successfully")
	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) StoreURL(ctx context.Context, shortCode, originalURL string) error {
	// Set with expiration (URLs expire after 1 year)
	expiration := 365 * 24 * time.Hour
	err := r.client.Set(ctx, shortCode, originalURL, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store URL in Redis: %w", err)
	}
	return nil
}

func (r *RedisStorage) GetURL(ctx context.Context, shortCode string) (string, error) {
	originalURL, err := r.client.Get(ctx, shortCode).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("URL not found")
		}
		return "", fmt.Errorf("failed to get URL from Redis: %w", err)
	}
	return originalURL, nil
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
