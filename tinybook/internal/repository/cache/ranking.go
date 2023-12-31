package cache

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"github.com/Yiling-J/theine-go"
	"github.com/bytedance/sonic"
	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, value []domain.Article, expiration time.Duration) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type RedisRankingCache struct {
	cli redis.Cmdable
	key string
}

func NewRedisRankingCache(cli redis.Cmdable) RankingCache {
	return &RedisRankingCache{cli: cli, key: "ranking:top_n"}
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	bytes, err := r.cli.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var articles []domain.Article
	err = sonic.Unmarshal(bytes, &articles)
	return articles, err
}

func (r *RedisRankingCache) Set(ctx context.Context, value []domain.Article, expiration time.Duration) error {
	for i := range value {
		value[i].Content = value[i].Abstract // 这里为了节省空间，将Content字段替换为Abstract字段
	}
	bytes, err := sonic.Marshal(value)
	if err != nil {
		return err
	}
	return r.cli.Set(ctx, r.key, bytes, expiration).Err()
}

type LocalRankingCache struct {
	cli *theine.Cache[string, any]
}

func NewLocalRankingCache(cli *theine.Cache[string, any]) RankingCache {
	return &LocalRankingCache{cli: cli}
}

func (l *LocalRankingCache) Set(ctx context.Context, value []domain.Article, expiration time.Duration) error {
	ttl := l.cli.SetWithTTL("ranking:top_n", value, 0, expiration)
	if ttl {
		return nil
	}
	return errors.New("local ranking本地缓存设置失败")
}

func (l *LocalRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	value, ok := l.cli.Get("ranking:top_n")
	if !ok {
		return nil, errors.New("local ranking本地缓存获取失败, 可能是缓存过期")
	}
	return value.([]domain.Article), nil
}
