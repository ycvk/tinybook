package service

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository"
	"github.com/samber/lo"
	"strconv"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	GetArticlesByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]domain.ArticleVo, error)
	GetArticleById(ctx context.Context, id int64) (domain.ArticleVo, error)
	GetPubArticleById(ctx context.Context, id int64) (domain.ArticleVo, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func (a *articleService) GetPubArticleById(ctx context.Context, id int64) (domain.ArticleVo, error) {
	article, err := a.repo.GetPubArticleById(ctx, id)
	if err != nil {
		return domain.ArticleVo{}, err
	}
	return domain.ArticleVo{
		ID:         article.ID,
		Title:      article.Title,
		Content:    article.Content,
		Author:     strconv.FormatInt(article.Author.ID, 10),
		AuthorName: article.Author.Name,
		Status:     strconv.FormatUint(uint64(article.Status), 10),
		Ctime:      time.Unix(article.Ctime, 0).Format("2006-01-02 15:04:05"),
		Utime:      time.Unix(article.Utime, 0).Format("2006-01-02 15:04:05"),
	}, nil
}

func (a *articleService) GetArticleById(ctx context.Context, id int64) (domain.ArticleVo, error) {
	article, err := a.repo.GetArticleById(ctx, id)
	if err != nil {
		return domain.ArticleVo{}, err
	}
	return domain.ArticleVo{
		ID:      article.ID,
		Title:   article.Title,
		Content: article.Content,
		Author:  strconv.FormatInt(article.Author.ID, 10),
		Status:  strconv.FormatUint(uint64(article.Status), 10),
		Ctime:   time.Unix(article.Ctime, 0).Format("2006-01-02 15:04:05"),
		Utime:   time.Unix(article.Utime, 0).Format("2006-01-02 15:04:05"),
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

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
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
