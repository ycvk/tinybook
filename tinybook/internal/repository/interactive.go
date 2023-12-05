package repository

import (
	"context"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error
	DecreaseLikeCount(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, cache: cache}
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
