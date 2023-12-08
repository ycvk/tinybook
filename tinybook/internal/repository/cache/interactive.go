package cache

import (
	"context"
	"fmt"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/events/interactive"
	"github.com/Yiling-J/theine-go"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strconv"
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
	GetTopNLike(ctx context.Context, biz string, num int64) ([]domain.Interactive, error)
	SetTopNLike(ctx context.Context, biz string, interactives []domain.Interactive) error
}

type RedisInteractiveCache struct {
	log           *zap.Logger
	cli           redis.Cmdable
	localCli      *theine.Cache[string, any]
	likeRankEvent interactive.LikeRankEventProducer
}

func (r *RedisInteractiveCache) SetTopNLike(ctx context.Context, biz string, interactives []domain.Interactive) error {
	key := r.key(biz, 0, LikeCountKey)
	pipeline := r.cli.Pipeline()
	for i := range interactives {
		pipeline.ZAdd(ctx, key, redis.Z{
			Score:  float64(interactives[i].LikeCount),
			Member: interactives[i].BizId,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func NewRedisInteractiveCache(cli redis.Cmdable, log *zap.Logger, cache *theine.Cache[string, any], event interactive.LikeRankEventProducer) InteractiveCache {
	return &RedisInteractiveCache{cli: cli, log: log, localCli: cache, likeRankEvent: event}
}

func (r *RedisInteractiveCache) GetTopNLike(ctx context.Context, biz string, num int64) ([]domain.Interactive, error) {
	key := r.key(biz, 0, LikeCountKey)
	// 从本地缓存中获取 topN 文章的点赞数与id
	localRank, ok := r.localCli.Get(key)
	if ok {
		r.log.Info("get topN like rank from local cache")
		interactiveLocalList := localRank.([]domain.Interactive)
		if len(interactiveLocalList) >= int(num) {
			return interactiveLocalList[:num], nil
		} else {
			return interactiveLocalList, nil
		}
	}
	// 从redis中获取 topN 文章的点赞数与id
	topNLike, err := r.cli.ZRevRangeWithScores(ctx, key, 0, num-1).Result()
	if err != nil {
		return nil, err
	}
	// redis.Z转换为 domain.Interactive
	interactivesMap := lo.Map(topNLike, func(item redis.Z, index int) domain.Interactive {
		s := item.Member.(string)
		id, _ := strconv.ParseInt(s, 10, 64)
		return domain.Interactive{
			BizId:     id,
			LikeCount: int64(item.Score),
		}
	})
	// 发送点赞排行榜事件 理论上有本地缓存存在，不会走到这里
	go func() {
		err2 := r.likeRankEvent.ProduceLikeRankEvent(interactive.LikeRankEvent{
			Change: true,
		})
		if err2 != nil {
			r.log.Error("produce like rank event failed", zap.Error(err2))
		}
	}()
	return interactivesMap, nil
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
	var interac domain.Interactive

	eg.Go(func() error {
		var er error
		interac.LikeCount, er = r.cli.SCard(ctx, r.key(biz, id, LikeCountKey)).Result()
		if er != nil {
			r.log.Warn("redis get likes error", zap.Error(er))
		}
		return er
	})

	eg.Go(func() error {
		var er error
		interac.CollectCount, er = r.cli.SCard(ctx, r.key(biz, id, CollectCountKey)).Result()
		if er != nil {
			r.log.Warn("redis get collects error", zap.Error(er))
		}
		return er
	})

	eg.Go(func() error {
		var er error
		interac.ReadCount, er = r.cli.HGet(ctx, r.key(biz, id, ReadCountKey), ReadCountKey).Int64()
		if er != nil {
			r.log.Warn("redis get read count error", zap.Error(er))
		}
		return er
	})

	return interac, eg.Wait()
}

func (r *RedisInteractiveCache) IncreaseCollectCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := r.key(biz, id, CollectCountKey)
	return r.cli.SAdd(ctx, key, uid).Err()
}

func (r *RedisInteractiveCache) IncreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := r.key(biz, id, LikeCountKey)
	// zincrby article:like_count 1 id
	_, err := r.cli.ZIncrBy(ctx, key, 1, strconv.FormatInt(id, 10)).Result()
	return err
}

func (r *RedisInteractiveCache) DecreaseLikeCountIfPresent(ctx context.Context, biz string, id int64, uid int64) error {
	key := r.key(biz, id, LikeCountKey)
	// zincrby article:like_count -1 id
	return r.cli.ZIncrBy(ctx, key, -1, strconv.FormatInt(id, 10)).Err()
}

func (r *RedisInteractiveCache) IncreaseReadCountIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.cli.HIncrBy(ctx, r.key(biz, bizId, ReadCountKey), ReadCountKey, 1).Err()
}

func (r *RedisInteractiveCache) key(biz string, bizId int64, keyType string) string {
	switch keyType {
	case ReadCountKey:
		return fmt.Sprintf("%s:%d", biz, bizId)
	case LikeCountKey:
		return fmt.Sprintf("%s:%s", biz, LikeCountKey)
	case CollectCountKey:
		return fmt.Sprintf("%s:%d:%s", biz, bizId, CollectCountKey)
	default:
		return fmt.Sprintf("%s:%d", biz, bizId)
	}
}
