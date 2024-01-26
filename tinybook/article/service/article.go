package service

import (
	"context"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strconv"
	"time"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/article/domain"
	"tinybook/tinybook/article/events/readcount"
	"tinybook/tinybook/article/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]domain.ArticleVo, error)
	GetArticleById(ctx context.Context, id int64) (domain.ArticleVo, error)
	GetPubArticleById(ctx context.Context, id int64, uid int64) (domain.ArticleVo, error)
	ListPub(ctx context.Context, time time.Time, limit int, offset int) ([]domain.Article, error)
	// 以下都是interactive service 的接口
	GetInteractive(ctx context.Context, request *intrv1.GetInteractiveRequest) (*intrv1.GetInteractiveResponse, error)
	Like(c context.Context, i *intrv1.LikeRequest) (*intrv1.LikeResponse, error)
	Unlike(c context.Context, i *intrv1.UnlikeRequest) (*intrv1.UnlikeResponse, error)
	Collect(ctx context.Context, i *intrv1.CollectRequest) (*intrv1.CollectResponse, error)
	GetLikeRanks(c context.Context, i *intrv1.GetLikeRanksRequest) (*intrv1.GetLikeRanksResponse, error)
	GetByIds(ctx context.Context, i *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error)
}

type articleService struct {
	repo     repository.ArticleRepository
	producer readcount.ReadEventProducer
	log      *zap.Logger
}

func NewArticleService(repo repository.ArticleRepository, producer readcount.ReadEventProducer, log *zap.Logger) ArticleService {
	//logger.With(zap.String("type", "articleService"))
	return &articleService{repo: repo, producer: producer, log: log}
}

func (a *articleService) GetByIds(ctx context.Context, i *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	return a.repo.GetByIds(ctx, i)
}

func (a *articleService) Unlike(c context.Context, i *intrv1.UnlikeRequest) (*intrv1.UnlikeResponse, error) {
	return a.repo.Unlike(c, i)
}

func (a *articleService) Collect(ctx context.Context, i *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	return a.repo.Collect(ctx, i)
}

func (a *articleService) GetLikeRanks(c context.Context, i *intrv1.GetLikeRanksRequest) (*intrv1.GetLikeRanksResponse, error) {
	return a.repo.GetLikeRanks(c, i)
}

func (a *articleService) Like(c context.Context, i *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	return a.repo.Like(c, i)
}

func (a *articleService) GetInteractive(ctx context.Context, request *intrv1.GetInteractiveRequest) (*intrv1.GetInteractiveResponse, error) {
	return a.repo.GetInteractive(ctx, request)
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
