package cache

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, value []domain.Article, expiration time.Duration) error
}

type RedisRankingCache struct {
	cli redis.Cmdable
	key string
}

func NewRedisRankingCache(cli redis.Cmdable) RankingCache {
	return &RedisRankingCache{cli: cli, key: "ranking:top_n"}
}

func (r *RedisRankingCache) Set(ctx context.Context, value []domain.Article, expiration time.Duration) error {
	for i := range value {
		value[i].Content = value[i].Abstract // 这里为了节省空间，将Content字段替换为Abstract字段
	}
	return r.cli.Set(ctx, r.key, value, expiration).Err()
}
