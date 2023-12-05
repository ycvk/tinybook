package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	ReadCountKey    = "read_count"
	LikeCountKey    = "like_count"
	CollectCountKey = "collect_count"
)

type InteractiveCache interface {
	IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error
	DecreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error
	IncreaseCollectCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error
}

type RedisInteractiveCache struct {
	cli redis.Cmdable
}

func NewRedisInteractiveCache(cli redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{cli: cli}
}

func (r *RedisInteractiveCache) IncreaseCollectCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := fmt.Sprintf("%s:%d:%s", biz, id, CollectCountKey)
	return r.cli.SAdd(ctx, key, uid).Err()
}

func (r *RedisInteractiveCache) IncreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := fmt.Sprintf("%s:%d:%s", biz, id, LikeCountKey)
	return r.cli.SAdd(ctx, key, uid).Err()
}

func (r *RedisInteractiveCache) DecreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := fmt.Sprintf("%s:%d:%s", biz, id, LikeCountKey)
	return r.cli.SRem(ctx, key, uid).Err()
}

func (r *RedisInteractiveCache) IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.cli.HIncrBy(ctx, r.key(biz, bizId), ReadCountKey, 1).Err()
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
