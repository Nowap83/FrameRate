package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	rdb *redis.Client
}

func NewCacheService(rdb *redis.Client) *CacheService {
	return &CacheService{rdb: rdb}
}

// récupère une donnée du cache et la dé-sérialise
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // miss
	} else if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}

	return true, nil // hit
}

// sérialise une donnée et la stocke dans le cache avec un TTL
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, key, data, expiration).Err()
}

// supprime une clé du cache
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.rdb.Del(ctx, key).Err()
}
