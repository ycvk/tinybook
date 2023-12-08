package repository

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"time"
)

type InteractiveRepository interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	BatchIncreaseReadCount(ctx context.Context, biz string, bizIds []int64) error
	IncreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error
	DecreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	GetInteractive(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetLikeRanks(ctx context.Context, biz string, num int64) ([]domain.Interactive, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	log   *zap.Logger
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache, logger *zap.Logger) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, cache: cache, log: logger}
}

func (c *CachedInteractiveRepository) GetLikeRanks(ctx context.Context, biz string, num int64) ([]domain.Interactive, error) {
	// 从缓存中获取 topN 文章的点赞数与id
	topNLikeCache, err2 := c.cache.GetTopNLike(ctx, biz, num)
	if err2 == nil && len(topNLikeCache) > 0 { // 缓存命中
		return topNLikeCache, nil
	}
	// 从数据库中获取 topN 文章的点赞数与id
	topNLike, err := c.dao.SelectTopNLike(ctx, biz, num)
	if err != nil {
		return nil, err
	}
	// dao转换为 domain.Interactive
	interactives := lo.Map(topNLike, func(item dao.Interactive, index int) domain.Interactive {
		return c.daoToDomain(item)
	})
	// 缓存
	go func() {
		timeout, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
		defer cancelFunc()
		err = c.cache.SetTopNLike(timeout, biz, interactives)
		if err != nil {
			c.log.Error("set topN like to cache failed", zap.Error(err))
		}
	}()
	return interactives, nil
}

func (c *CachedInteractiveRepository) BatchIncreaseReadCount(ctx context.Context, biz string, bizIds []int64) error {
	err := c.dao.BatchIncreaseReadCount(ctx, biz, bizIds)
	if err != nil {
		return err
	}
	return c.cache.BatchIncreaseReadCountIfPresent(ctx, biz, bizIds)
}

func (c *CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	collected, err := c.cache.IsCollected(ctx, biz, id, uid)
	if err == nil {
		return collected, nil
	}
	return c.dao.IsCollected(ctx, biz, id, uid)
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	liked, err := c.cache.IsLiked(ctx, biz, id, uid)
	if err == nil {
		return liked, nil
	}
	return c.dao.IsLiked(ctx, biz, id, uid)
}

func (c *CachedInteractiveRepository) GetInteractive(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	interactive, err := c.cache.GetInteractive(ctx, biz, id)
	if err == nil {
		return interactive, nil
	}
	getInteractive, err := c.dao.GetInteractive(ctx, biz, id)
	if err != nil {
		if err.Error() == dao.ErrNotFound {
			return domain.Interactive{}, nil
		}
		return domain.Interactive{}, err
	}
	interactive = c.daoToDomain(getInteractive)
	return interactive, nil
}

func (c *CachedInteractiveRepository) Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	err := c.dao.InsertCollectRecord(ctx, biz, id, cid, uid)
	if err != nil {
		return err
	}
	return c.cache.IncreaseCollectCountIfPresent(ctx, biz, id, uid)
}

func (c *CachedInteractiveRepository) IncreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.InsertLikeRecord(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.IncreaseLikeCountIfPresent(ctx, biz, id, uid)
}

func (c *CachedInteractiveRepository) DecreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.DeleteLikeRecord(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.DecreaseLikeCountIfPresent(ctx, biz, id, uid)
}

func (c *CachedInteractiveRepository) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncreaseReadCount(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return c.cache.IncreaseReadCountIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) daoToDomain(interactive dao.Interactive) domain.Interactive {
	return domain.Interactive{
		Biz:   interactive.Biz,
		BizId: interactive.BizId,

		LikeCount:    interactive.LikeCount,
		CollectCount: interactive.CollectCount,
		ReadCount:    interactive.ReadCount,
	}
}
