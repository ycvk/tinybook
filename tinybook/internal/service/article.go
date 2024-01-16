package service

import (
	"context"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strconv"
	"time"
	"tinybook/tinybook/internal/domain"
	"tinybook/tinybook/internal/events/readcount"
	"tinybook/tinybook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]domain.ArticleVo, error)
	GetArticleById(ctx context.Context, id int64) (domain.ArticleVo, error)
	GetPubArticleById(ctx context.Context, id int64, uid int64) (domain.ArticleVo, error)
	ListPub(ctx context.Context, time time.Time, limit int, offset int) ([]domain.Article, error)
}

type articleService struct {
	repo     repository.ArticleRepository
	producer readcount.ReadEventProducer
	log      *zap.Logger
}

func (a *articleService) ListPub(ctx context.Context, time time.Time, limit int, offset int) ([]domain.Article, error) {
	return a.repo.ListPub(ctx, time, limit, offset)
}

func (a *articleService) GetPubArticleById(ctx context.Context, id int64, uid int64) (domain.ArticleVo, error) {
	art, err := a.repo.GetPubArticleById(ctx, id)
	if err != nil {
		return domain.ArticleVo{}, err
	}
	// 异步发送阅读事件
	go func() {
		err = a.producer.ProduceReadEvent(readcount.ReadEvent{
			ArticleID: id,
			UserID:    uid,
		})
		if err != nil {
			a.log.Error("produce read event failed, article id: "+
				strconv.FormatInt(id, 10)+" user id: "+
				strconv.FormatInt(uid, 10), zap.Error(err))
		}
	}()
	return domain.ArticleVo{
		ID:         art.ID,
		Title:      art.Title,
		Content:    art.Content,
		Author:     strconv.FormatInt(art.Author.ID, 10),
		AuthorName: art.Author.Name,
		Status:     strconv.FormatUint(uint64(art.Status), 10),
		Ctime:      time.Unix(art.Ctime, 0).Format("2006-01-02 15:04:05"),
		Utime:      time.Unix(art.Utime, 0).Format("2006-01-02 15:04:05"),
	}, nil
}

func (a *articleService) GetArticleById(ctx context.Context, id int64) (domain.ArticleVo, error) {
	art, err := a.repo.GetArticleById(ctx, id)
	if err != nil {
		return domain.ArticleVo{}, err
	}
	return domain.ArticleVo{
		ID:      art.ID,
		Title:   art.Title,
		Content: art.Content,
		Author:  strconv.FormatInt(art.Author.ID, 10),
		Status:  strconv.FormatUint(uint64(art.Status), 10),
		Ctime:   time.Unix(art.Ctime, 0).Format("2006-01-02 15:04:05"),
		Utime:   time.Unix(art.Utime, 0).Format("2006-01-02 15:04:05"),
	}, nil
}

func (a *articleService) GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]domain.ArticleVo, error) {
	articles, err := a.repo.GetArticlesByAuthor(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	return lo.Map(articles, func(arts domain.Article, index int) domain.ArticleVo {
		return domain.ArticleVo{
			ID:    arts.ID,
			Title: arts.Title,
			//Content:  arts.Content,
			Abstract: arts.Abstract,
			Author:   strconv.FormatInt(arts.Author.ID, 10),
			Status:   strconv.FormatUint(uint64(arts.Status), 10),
			Ctime:    time.Unix(arts.Ctime, 0).Format("2006-01-02 15:04:05"),
			Utime:    time.Unix(arts.Utime, 0).Format("2006-01-02 15:04:05"),
		}
	}), nil
}

func NewArticleService(repo repository.ArticleRepository, producer readcount.ReadEventProducer, log *zap.Logger) ArticleService {
	//logger.With(zap.String("type", "articleService"))
	return &articleService{repo: repo, producer: producer, log: log}
}

func (a *articleService) Withdraw(ctx context.Context, article domain.Article) error {
	return a.repo.SyncStatus(ctx, article, domain.ArticleStatusPrivate) // 私有
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished // 已发布
	return a.repo.Sync(ctx, article)
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnpublished // 未发布
	if article.ID > 0 {
		return article.ID, a.repo.Update(ctx, article)
	}
	return a.repo.Create(ctx, article)
}
