package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, duration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisArticleCache struct {
	cli redis.Cmdable
}

func (r *RedisArticleCache) Delete(ctx context.Context, key string) error {
	return r.cli.Del(ctx, key).Err()
}

func (r *RedisArticleCache) Set(ctx context.Context, key string, value []byte, duration time.Duration) error {
	return r.cli.Set(ctx, key, value, duration).Err()
}

func (r *RedisArticleCache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.cli.Get(ctx, key).Bytes()
}

func NewRedisArticleCache(cli redis.Cmdable) ArticleCache {
	return &RedisArticleCache{cli: cli}
}
