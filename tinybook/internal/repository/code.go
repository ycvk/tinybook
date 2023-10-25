package repository

import (
	"context"
	"geek_homework/tinybook/internal/repository/cache"
)

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code, timeInterval string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CachedCodeRepository struct {
	cache cache.CodeCache
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
