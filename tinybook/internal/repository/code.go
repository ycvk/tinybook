package repository

import (
	"context"
	"geek_homework/tinybook/internal/repository/cache"
)

type CodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: *cache,
	}
}

func (repo *CodeRepository) Set(ctx context.Context, biz, phone, code, timeInterval string) error {
	return repo.cache.SetCode(ctx, phone, biz, code, timeInterval)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.VerifyCode(ctx, phone, biz, code)
}
