package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const ReadCountKey = "read_count"

type InteractiveCache interface {
	IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error
}

type RedisInteractiveCache struct {
	cli redis.Cmdable
}

func NewRedisInteractiveCache(cli redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{cli: cli}
}

func (r *RedisInteractiveCache) IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.cli.HIncrBy(ctx, r.key(biz, bizId), ReadCountKey, 1).Err()
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
