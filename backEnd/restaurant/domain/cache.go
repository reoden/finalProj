package domain

import (
	"context"
	"restaurant/pkgs"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	db *redis.Client
}

func NewCacheRepo() *CacheRepo {
	return &CacheRepo{
		db: pkgs.GetRedis(),
	}
}

func (cr *CacheRepo) SetKey(ctx context.Context, key, value string, seconds time.Duration) error {
	return cr.db.Set(ctx, key, value, seconds).Err()
}

func (cr *CacheRepo) IncrKey(ctx context.Context, key string) error {
	return cr.db.Incr(ctx, key).Err()
}

func (cr *CacheRepo) DeleteKey(ctx context.Context, key string) error {
	return cr.db.Del(ctx, key).Err()
}

func (cr *CacheRepo) TTLKey(ctx context.Context, key string) (int64, error) {
	d, err := cr.db.TTL(ctx, key).Result()
	return int64(d.Seconds()), err
}

func (cr *CacheRepo) GetKey(ctx context.Context, key string) (string, error) {
	value, err := cr.db.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (cr *CacheRepo) LPush(ctx context.Context, key string, data []byte) (int64, error) {
	return cr.db.LPush(context.Background(), key, data).Result()
}
