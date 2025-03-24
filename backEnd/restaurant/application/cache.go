package application

import (
	"context"
	"log"
	"restaurant/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	db *domain.CacheRepo
}

func NewCacheService() *CacheService {
	return &CacheService{
		db: domain.NewCacheRepo(),
	}
}

func (c *CacheService) SetKey(ctx context.Context, key, value string, seconds time.Duration) error {
	return c.db.SetKey(ctx, key, value, seconds)
}

func (c *CacheService) IncrKey(ctx context.Context, key string) error {
	return c.db.IncrKey(ctx, key)
}

func (c *CacheService) DeleteKey(ctx context.Context, key string) error {
	return c.db.DeleteKey(ctx, key)
}

func (c *CacheService) TTLKey(ctx context.Context, key string) (int64, error) {
	return c.db.TTLKey(ctx, key)
}

func (c *CacheService) GetKey(ctx context.Context, key string) (string, error) {
	data, err := c.db.GetKey(ctx, key)
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return data, nil
}

func (c *CacheService) PushQueue(ctx context.Context, key string, data []byte) error {
	_, err := c.db.LPush(ctx, key, data)
	return err
}

func (a *CacheService) setCodeCache(ctx context.Context, email, code string) error {
	return a.SetKey(ctx, email, code, 60*time.Second*3)
}

func (a *CacheService) getCodeCache(ctx context.Context, email string) (string, error) {
	data, err := a.GetKey(ctx, email)
	if err != nil {
		log.Printf("getCodeCache_GetKey_err: %+v", err)
		return "", err
	}

	return data, nil
}
