package repository

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error
	DecreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	GetInteractive(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, cache: cache}
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
		LikeCount:    interactive.LikeCount,
		CollectCount: interactive.CollectCount,
		ReadCount:    interactive.ReadCount,
	}
}
