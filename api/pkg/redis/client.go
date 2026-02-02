package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var client *redis.Client

func Init(addr, password string, db int) error {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}

func GetClient() *redis.Client {
	return client
}

func Get(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

func Del(ctx context.Context, keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

func Incr(ctx context.Context, key string) error {
	return client.Incr(ctx, key).Err()
}

func HSet(ctx context.Context, key string, values ...interface{}) error {
	return client.HSet(ctx, key, values...).Err()
}

func HGet(ctx context.Context, key, field string) (string, error) {
	return client.HGet(ctx, key, field).Result()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return client.HGetAll(ctx, key).Result()
}

func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return client.ZAdd(ctx, key, members...).Err()
}

func ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func PFAdd(ctx context.Context, key string, els ...interface{}) error {
	return client.PFAdd(ctx, key, els...).Err()
}

func PFCount(ctx context.Context, keys ...string) (int64, error) {
	return client.PFCount(ctx, keys...).Result()
}
