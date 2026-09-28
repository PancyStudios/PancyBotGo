package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore is a shared cache backed by Redis.
type RedisStore struct {
	client *redis.Client
	prefix string
}

// NewRedisStore connects to Redis using a redis:// or rediss:// URL.
func NewRedisStore(ctx context.Context, rawURL, prefix string) (*RedisStore, error) {
	prefix = strings.TrimSuffix(strings.TrimSpace(prefix), ":")
	if prefix == "" {
		return nil, errors.New("REDIS_PREFIX cannot be empty")
	}
	for _, character := range prefix {
		if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789:._-", character) {
			return nil, errors.New("REDIS_PREFIX contains unsupported characters")
		}
	}
	prefix += ":"

	options, err := redis.ParseURL(rawURL)
	if err != nil {
		return nil, errors.New("invalid REDIS_URL")
	}
	options.DialTimeout = 3 * time.Second
	options.ReadTimeout = 750 * time.Millisecond
	options.WriteTimeout = 750 * time.Millisecond

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping Redis: %w", err)
	}

	return &RedisStore{client: client, prefix: prefix}, nil
}

func (r *RedisStore) key(key string) string { return r.prefix + key }

func (r *RedisStore) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := r.client.Get(ctx, r.key(key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrMiss
	}
	return value, err
}

func (r *RedisStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, r.key(key), value, ttl).Err()
}

func (r *RedisStore) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.key(key)).Err()
}

func (r *RedisStore) Clear(ctx context.Context) error {
	var cursor uint64
	for {
		keys, next, err := r.client.Scan(ctx, cursor, r.prefix+"*", 250).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := r.client.Unlink(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func (r *RedisStore) Size(ctx context.Context) (int, error) {
	keys, err := r.keys(ctx)
	return len(keys), err
}

func (r *RedisStore) keys(ctx context.Context) ([]string, error) {
	pattern := r.prefix + "*"
	var cursor uint64
	keys := make([]string, 0)
	for {
		batch, next, err := r.client.Scan(ctx, cursor, pattern, 250).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, batch...)
		cursor = next
		if cursor == 0 {
			return keys, nil
		}
	}
}

func (r *RedisStore) Close() error { return r.client.Close() }
func (r *RedisStore) Name() string { return "redis" }
