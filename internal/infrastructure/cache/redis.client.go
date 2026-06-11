package cache

import (
	"go-platform-core/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
