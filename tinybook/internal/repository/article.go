package repository

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, article domain.Article, articleStatus domain.ArticleStatus) error
	GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	GetFirstPage(ctx context.Context, uid int64, limit int) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
}

type CachedArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
	log   *zap.Logger
}

func (c *CachedArticleRepository) DelFirstPage(ctx context.Context, uid int64) error {
	return c.cache.Delete(ctx, c.GetCacheKey(uid))
}

func (c *CachedArticleRepository) SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error {
	key := c.GetCacheKey(uid)
	// 只需要缓存文章的摘要
	i := lo.Map(articles, func(article domain.Article, index int) domain.Article {
		article.Content = article.Abstract
		return article
	})
	marshal, err := sonic.Marshal(i)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, marshal, 30*time.Minute)
}

func (c *CachedArticleRepository) GetCacheKey(uid int64) string {
	return strconv.FormatInt(uid, 10) + "_first_page"
}

func (c *CachedArticleRepository) GetFirstPage(ctx context.Context, uid int64, limit int) ([]domain.Article, error) {
	bytes, err := c.cache.Get(ctx, c.GetCacheKey(uid))
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = sonic.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}
	if len(res) < limit {
		return res, nil
	}
	return res[:limit], nil
}

func (c *CachedArticleRepository) GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 { //如果是第一页，且limit小于100，从缓存中取
		firstPage, err := c.GetFirstPage(ctx, uid, limit)
		if err == nil {
			return firstPage, nil
		}
		c.log.Error("get first page from cache failed", zap.Error(err))
	}
	articles, err := c.dao.GetArticlesByAuthor(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	go func() { //异步更新缓存
		firstPageRes, err2 := c.dao.GetArticlesByAuthor(ctx, uid, 100, 0)
		if err2 != nil {
			c.log.Error("get first page from db failed", zap.Error(err2))
			return
		}
		domainRes := lo.Map(firstPageRes, func(article dao.Article, index int) domain.Article {
			return c.daoToDomain(article)
		})
		err2 = c.SetFirstPage(ctx, uid, domainRes)
		if err2 != nil {
			c.log.Error("set first page to cache failed", zap.Error(err2))
			return
		}
	}()
	return lo.Map(articles, func(article dao.Article, index int) domain.Article {
		return c.daoToDomain(article)
	}), nil
}

func NewCachedArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, log *zap.Logger) ArticleRepository {
	return &CachedArticleRepository{dao: dao, cache: cache, log: log}
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, article domain.Article, articleStatus domain.ArticleStatus) error {
	err := c.dao.SyncStatus(ctx, c.domainToDao(article), uint8(articleStatus))
	delErr := c.DelFirstPage(ctx, article.Author.ID)
	if delErr != nil {
		c.log.Error("delete first page from cache failed", zap.Error(delErr))
	}
	return err
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	sync, err := c.dao.Sync(ctx, c.domainToDao(article))
	delErr := c.DelFirstPage(ctx, article.Author.ID)
	if delErr != nil {
		c.log.Error("delete first page from cache failed", zap.Error(delErr))
	}
	return sync, err
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	err := c.dao.UpdateById(ctx, c.domainToDao(article))
	delErr := c.DelFirstPage(ctx, article.Author.ID)
	if delErr != nil {
		c.log.Error("delete first page from cache failed", zap.Error(delErr))
	}
	return err
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	insert, err := c.dao.Insert(ctx, c.domainToDao(article))
	delErr := c.DelFirstPage(ctx, article.Author.ID)
	if delErr != nil {
		c.log.Error("delete first page from cache failed", zap.Error(delErr))
	}
	return insert, err
}

func (c *CachedArticleRepository) domainToDao(article domain.Article) dao.Article {
	return dao.Article{
		ID:       article.ID,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.ID,
		Status:   uint8(article.Status),
	}
}

func (c *CachedArticleRepository) daoToDomain(article dao.Article) domain.Article {
	var abstract string
	runeContent := []rune(article.Content)
	// 取前128个字符作为摘要
	if len(runeContent) > 128 {
		abstract = string(runeContent[:128])
	} else {
		abstract = string(runeContent)
	}
	return domain.Article{
		ID:       article.ID,
		Title:    article.Title,
		Content:  article.Content,
		Abstract: abstract,
		Author:   domain.Author{ID: article.AuthorId},
		Status:   domain.ArticleStatus(article.Status),
		Ctime:    article.Ctime,
		Utime:    article.Utime,
	}
}
