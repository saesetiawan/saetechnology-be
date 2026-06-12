package cache

import (
	"context"
	"time"

	cacheDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/cache"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RedisRepository struct {
	Client *redis.Client
	Tracer trace.Tracer
}

func NewRedisRepository(
	client *redis.Client,
	tracerProvider trace.TracerProvider,
) cacheDomain.Repository {
	return &RedisRepository{
		Client: client,
		Tracer: tracerProvider.Tracer("RedisRepository"),
	}
}

func (r *RedisRepository) Get(
	ctx context.Context,
	key string,
) (string, error) {
	ctx, span := r.Tracer.Start(ctx, "RedisRepository.Get")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.provider", "redis"),
		attribute.String("cache.operation", "get"),
		attribute.String("cache.key", key),
	)

	result, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			span.SetAttributes(
				attribute.Bool("cache.hit", false),
			)

			span.SetStatus(codes.Ok, "cache miss")

			return "", err
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, "failed get cache")

		return "", err
	}

	span.SetAttributes(
		attribute.Bool("cache.hit", true),
		attribute.Int("cache.value_length", len(result)),
	)

	span.SetStatus(codes.Ok, "cache hit")

	return result, nil
}

func (r *RedisRepository) Set(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {
	ctx, span := r.Tracer.Start(ctx, "RedisRepository.Set")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.provider", "redis"),
		attribute.String("cache.operation", "set"),
		attribute.String("cache.key", key),
		attribute.Int("cache.value_length", len(value)),
		attribute.Int64("cache.ttl_seconds", int64(ttl.Seconds())),
	)

	err := r.Client.Set(
		ctx,
		key,
		value,
		ttl,
	).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed set cache")
		return err
	}

	span.SetStatus(codes.Ok, "cache saved")

	return nil
}

func (r *RedisRepository) Delete(
	ctx context.Context,
	key string,
) error {
	ctx, span := r.Tracer.Start(ctx, "RedisRepository.Delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("cache.provider", "redis"),
		attribute.String("cache.operation", "delete"),
		attribute.String("cache.key", key),
	)

	result, err := r.Client.Del(ctx, key).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed delete cache")
		return err
	}

	span.SetAttributes(
		attribute.Int64("cache.deleted_count", result),
	)

	span.SetStatus(codes.Ok, "cache deleted")

	return nil
}

func (r *RedisRepository) DeleteByPrefix(
	ctx context.Context,
	prefix string,
) error {
	ctx, span := r.Tracer.Start(ctx, "RedisRepository.DeleteByPrefix")
	defer span.End()

	pattern := prefix + "*"

	span.SetAttributes(
		attribute.String("cache.provider", "redis"),
		attribute.String("cache.operation", "delete_by_prefix"),
		attribute.String("cache.prefix", prefix),
		attribute.String("cache.pattern", pattern),
	)

	var cursor uint64
	var deletedCount int64

	for {
		keys, nextCursor, err := r.Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed scan cache keys")
			return err
		}

		if len(keys) > 0 {
			result, err := r.Client.Del(ctx, keys...).Result()
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed delete cache keys")
				return err
			}

			deletedCount += result
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	span.SetAttributes(
		attribute.Int64("cache.deleted_count", deletedCount),
	)

	span.SetStatus(codes.Ok, "cache deleted by prefix")

	return nil
}
