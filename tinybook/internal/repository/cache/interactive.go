package cache

import (
	"context"
	"fmt"
	"geek_homework/tinybook/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	ReadCountKey    = "read_count"
	LikeCountKey    = "like_count"
	CollectCountKey = "collect_count"
)

type InteractiveCache interface {
	IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error
	BatchIncreaseReadCountIfPresent(ctx context.Context, biz string, ids []int64) error
	IncreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error
	DecreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error
	IncreaseCollectCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error
	GetInteractive(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	IsLiked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	IsCollected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type RedisInteractiveCache struct {
	log *zap.Logger
	cli redis.Cmdable
}

func NewRedisInteractiveCache(cli redis.Cmdable, log *zap.Logger) InteractiveCache {
	return &RedisInteractiveCache{cli: cli, log: log}
}

func (r *RedisInteractiveCache) BatchIncreaseReadCountIfPresent(ctx context.Context, biz string, ids []int64) error {
	var eg errgroup.Group
	for i := range ids {
		id := ids[i]
		eg.Go(func() error {
			return r.IncreaseReadCountIfPresent(ctx, biz, id)
		})
	}
	return eg.Wait()
}

func (r *RedisInteractiveCache) IsCollected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	return r.cli.SIsMember(ctx, r.key(biz, id, CollectCountKey), uid).Result()
}

func (r *RedisInteractiveCache) IsLiked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	return r.cli.SIsMember(ctx, r.key(biz, id, LikeCountKey), uid).Result()
}

func (r *RedisInteractiveCache) GetInteractive(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	var eg errgroup.Group
	var interactive domain.Interactive

	eg.Go(func() error {
		var er error
		interactive.LikeCount, er = r.cli.SCard(ctx, r.key(biz, id, LikeCountKey)).Result()
		if er != nil {
			r.log.Warn("redis get likes error", zap.Error(er))
		}
		return er
	})

	eg.Go(func() error {
		var er error
		interactive.CollectCount, er = r.cli.SCard(ctx, r.key(biz, id, CollectCountKey)).Result()
		if er != nil {
			r.log.Warn("redis get collects error", zap.Error(er))
		}
		return er
	})

	eg.Go(func() error {
		var er error
		interactive.ReadCount, er = r.cli.HGet(ctx, r.key(biz, id, ReadCountKey), ReadCountKey).Int64()
		if er != nil {
			r.log.Warn("redis get read count error", zap.Error(er))
		}
		return er
	})

	return interactive, eg.Wait()
}

func (r *RedisInteractiveCache) IncreaseCollectCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := r.key(biz, id, CollectCountKey)
	return r.cli.SAdd(ctx, key, uid).Err()
}

func (r *RedisInteractiveCache) IncreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := r.key(biz, id, LikeCountKey)
	return r.cli.SAdd(ctx, key, uid).Err()
}

func (r *RedisInteractiveCache) DecreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := r.key(biz, id, LikeCountKey)
	return r.cli.SRem(ctx, key, uid).Err()
}

func (r *RedisInteractiveCache) IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.cli.HIncrBy(ctx, r.key(biz, bizId, ReadCountKey), ReadCountKey, 1).Err()
}

func (r *RedisInteractiveCache) key(biz string, bizId int64, keyType string) string {
	switch keyType {
	case ReadCountKey:
		return fmt.Sprintf("%s:%d", biz, bizId)
	case LikeCountKey:
		return fmt.Sprintf("%s:%d:%s", biz, bizId, LikeCountKey)
	case CollectCountKey:
		return fmt.Sprintf("%s:%d:%s", biz, bizId, CollectCountKey)
	default:
		return fmt.Sprintf("%s:%d", biz, bizId)
	}
}
