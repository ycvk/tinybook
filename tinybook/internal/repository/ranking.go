package repository

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, topN []domain.Article) error
}

type CachedRankingRepository struct {
	cache cache.RankingCache
}

func NewCachedRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CachedRankingRepository{cache: cache}
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, topN []domain.Article) error {
	return c.cache.Set(ctx, topN, 0)
}
