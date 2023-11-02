package repository

import (
	"context"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
)

var DuplicatePhoneError = dao.DuplicatePhoneError

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code, timeInterval string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
	Create(ctx context.Context, number ...string) error
	Delete(ctx context.Context, number string) error
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

type GormCodeRepository struct {
	dao dao.CodeDAO
}

func (g *GormCodeRepository) Create(ctx context.Context, number ...string) error {
	codes := make([]*dao.Code, len(number))
	for _, num := range number {
		codes = append(codes, &dao.Code{
			Number: num,
		})
	}
	return g.dao.Insert(ctx, codes...)
}

func (g *GormCodeRepository) Delete(ctx context.Context, number string) error {
	return g.dao.Delete(ctx, dao.Code{
		Number: number,
	})
}

func NewCachedCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

func (repo *CachedCodeRepository) Set(ctx context.Context, biz, phone, code, timeInterval string) error {
	return repo.cache.SetCode(ctx, phone, biz, code, timeInterval)
}

func (repo *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.VerifyCode(ctx, phone, biz, code)
}

func (repo *CachedCodeRepository) Create(ctx context.Context, number ...string) error {
	panic("implement me")
}

func (repo *CachedCodeRepository) Delete(ctx context.Context, number string) error {
	panic("implement me")
}

func (g *GormCodeRepository) Set(ctx context.Context, biz, phone, code, timeInterval string) error {
	panic("implement me")
}

func (g *GormCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	panic("implement me")
}
