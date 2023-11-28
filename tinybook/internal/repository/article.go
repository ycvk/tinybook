package repository

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"github.com/bytedance/sonic"
	"github.com/cockroachdb/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	PubArticleKey       = "article:pub:"
	ArticleKey          = "article:"
	ArticleFirstPageKey = "article:first_page:"
)

type ArticleType int

const (
	ArticleUnknown ArticleType = iota
	ArticleAuthor
	ArticleReader
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
	GetArticleById(ctx context.Context, id int64) (domain.Article, error)
	SetCache(ctx context.Context, key int64, articleType ArticleType, article domain.Article, expire time.Duration) error
	GetCache(ctx context.Context, key int64, articleType ArticleType) (domain.Article, error)
	DelCache(ctx context.Context, key int64, articleType ArticleType) error
	GetPubArticleById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao      dao.ArticleDAO
	cache    cache.ArticleCache
	userRepo UserRepository
	log      *zap.Logger
}

func (c *CachedArticleRepository) GetPubArticleById(ctx context.Context, id int64) (domain.Article, error) {
	getCache, err := c.GetCache(ctx, id, ArticleReader)
	if err == nil {
		return getCache, nil
	}
	article, err := c.dao.GetPubArticleById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 获取作者信息
	user, err := c.userRepo.FindById(ctx, article.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	toDomain := c.daoToDomain(dao.Article(article))
	toDomain.Author.Name = user.Nickname
	go func() {
		timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelFunc()
		err2 := c.SetCache(timeout, article.ID, ArticleReader, toDomain, 3*24*time.Hour)
		if err2 != nil {
			c.log.Warn("preset published article to cache failed", zap.Error(err2))
			return
		}
	}()
	return toDomain, nil
}

func (c *CachedArticleRepository) SetCache(ctx context.Context, key int64, articleType ArticleType, art domain.Article, expire time.Duration) error {
	var articleKey string
	// 根据文章类型，选择不同的缓存key
	switch articleType {
	case ArticleAuthor:
		articleKey = c.GetCacheArticleKey(key)
	case ArticleReader:
		articleKey = c.GetPubCacheArticleKey(key)
	default:
		return errors.New("unknown article type")
	}
	marshal, err := sonic.Marshal(art)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, articleKey, marshal, expire)
}

func (c *CachedArticleRepository) GetCache(ctx context.Context, key int64, articleType ArticleType) (domain.Article, error) {
	var articleKey string
	// 根据文章类型，选择不同的缓存key
	switch articleType {
	case ArticleAuthor:
		articleKey = c.GetCacheArticleKey(key)
	case ArticleReader:
		articleKey = c.GetPubCacheArticleKey(key)
	default:
		return domain.Article{}, errors.New("unknown article type")
	}
	bytes, err := c.cache.Get(ctx, articleKey)
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = sonic.Unmarshal(bytes, &res)
	if err != nil {
		return domain.Article{}, err
	}
	return res, nil
}

func (c *CachedArticleRepository) DelCache(ctx context.Context, key int64, articleType ArticleType) error {
	var articleKey string
	// 根据文章类型，选择不同的缓存key
	switch articleType {
	case ArticleAuthor:
		articleKey = c.GetCacheArticleKey(key)
	case ArticleReader:
		articleKey = c.GetPubCacheArticleKey(key)
	default:
		return errors.New("unknown article type")
	}
	return c.cache.Delete(ctx, articleKey)
}

func (c *CachedArticleRepository) GetArticleById(ctx context.Context, id int64) (domain.Article, error) {
	cacheRes, err := c.GetCache(ctx, id, ArticleAuthor)
	if err == nil {
		return cacheRes, nil
	}
	article, err := c.dao.GetArticleById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.daoToDomain(article), nil
}

func (c *CachedArticleRepository) DelFirstPage(ctx context.Context, uid int64) error {
	return c.cache.Delete(ctx, c.GetCachePageKey(uid))
}

func (c *CachedArticleRepository) SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error {
	key := c.GetCachePageKey(uid)
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

func (c *CachedArticleRepository) GetCachePageKey(uid int64) string {
	return ArticleFirstPageKey + strconv.FormatInt(uid, 10)
}

func (c *CachedArticleRepository) GetCacheArticleKey(id int64) string {
	return ArticleKey + strconv.FormatInt(id, 10)
}

func (c *CachedArticleRepository) GetPubCacheArticleKey(id int64) string {
	return PubArticleKey + strconv.FormatInt(id, 10)
}

func (c *CachedArticleRepository) GetFirstPage(ctx context.Context, uid int64, limit int) ([]domain.Article, error) {
	bytes, err := c.cache.Get(ctx, c.GetCachePageKey(uid))
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
		if offset != 0 || limit > 100 { //如果不是第一页，或者limit大于100，没有必要更新缓存
			return
		}
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
	go func() { //异步预加载第一个文章到缓存
		if len(articles) > 0 {
			article := articles[0]
			timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancelFunc()
			err2 := c.SetCache(timeout, article.ID, ArticleAuthor, c.daoToDomain(article), time.Minute)
			if err2 != nil {
				c.log.Warn("preset article to cache failed", zap.Error(err2))
				return
			}
		}
	}()
	return lo.Map(articles, func(article dao.Article, index int) domain.Article {
		return c.daoToDomain(article)
	}), nil
}

func NewCachedArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, userRepo UserRepository, log *zap.Logger) ArticleRepository {
	return &CachedArticleRepository{dao: dao, cache: cache, userRepo: userRepo, log: log}
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, article domain.Article, articleStatus domain.ArticleStatus) error {
	err := c.dao.SyncStatus(ctx, c.domainToDao(article), uint8(articleStatus))
	delErr := c.DelFirstPage(ctx, article.Author.ID)
	// 删除读者缓存
	pubErr := c.DelCache(ctx, article.ID, ArticleReader)
	if delErr != nil {
		c.log.Error("delete first page from cache failed", zap.Error(delErr))
	}
	if pubErr != nil {
		c.log.Warn("delete article from cache failed", zap.Error(err))
	}
	return err
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	sync, err := c.dao.Sync(ctx, c.domainToDao(article))
	delErr := c.DelFirstPage(ctx, article.Author.ID)
	// 删除读者缓存
	pubErr := c.DelCache(ctx, article.ID, ArticleReader)
	if delErr != nil {
		c.log.Error("delete first page from cache failed", zap.Error(delErr))
	}
	if pubErr != nil {
		c.log.Warn("delete article from cache failed", zap.Error(err))
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
