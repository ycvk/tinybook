package service

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/events/interactive"
	"geek_homework/tinybook/internal/repository"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strconv"
)

type InteractiveService interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	Unlike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	GetInteractive(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	GetLikeRanks(ctx context.Context, biz string, num int64) ([]domain.ArticleVo, error)
}

type interactiveService struct {
	repo          repository.InteractiveRepository
	articleRepo   repository.ArticleRepository
	likeRankEvent interactive.LikeRankEventProducer
	log           *zap.Logger
}

func (i *interactiveService) GetLikeRanks(ctx context.Context, biz string, num int64) ([]domain.ArticleVo, error) {
	// 获取 topN 文章的点赞数与id
	likeRanks, err := i.repo.GetLikeRanks(ctx, biz, num)
	if err != nil {
		return nil, err
	}
	// 获取文章详情
	var eg errgroup.Group
	articles := make([]domain.Article, len(likeRanks))
	for i2 := range likeRanks {
		index := i2 // 这里需要注意，闭包问题
		eg.Go(func() error {
			// 获取单个文章详情
			article, err := i.articleRepo.GetPubArticleById(ctx, likeRanks[index].BizId) // `GetPubArticleById`已经走了缓存，所以这里不用再走缓存了
			if err != nil {
				return err
			}
			articles[index] = article // 这样赋值，可以保证与`likeRanks`的顺序一致
			return nil
		})
	}
	err = eg.Wait()
	if err != nil {
		return nil, err
	}
	// 转换为 ArticleVo
	vos := lo.Map(articles, func(item domain.Article, index int) domain.ArticleVo {
		return domain.ArticleVo{
			ID:          item.ID,
			Title:       item.Title,
			Content:     item.Content, // 这里返回的文章内容，是已经在`GetPubArticleById`里经过处理的，比如截取前100个字符
			Interactive: likeRanks[index],
		}
	})
	return vos, nil
}

func (i *interactiveService) GetInteractive(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	interactiveData, err := i.repo.GetInteractive(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		interactiveData.Liked, er = i.repo.Liked(ctx, biz, id, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		interactiveData.Collected, er = i.repo.Collected(ctx, biz, id, uid)
		return er
	})
	return interactiveData, eg.Wait()
}

func (i *interactiveService) Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return i.repo.Collect(ctx, biz, id, cid, uid)
}

func (i *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	err := i.repo.IncreaseLikeCount(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	// 异步发送点赞事件
	go func() {
		err2 := i.likeRankEvent.ProduceLikeRankEvent(interactive.LikeRankEvent{
			ArticleID: id,
			Change:    true,
		})
		if err2 != nil {
			i.log.Error("produce like rank event failed, article id: "+
				strconv.FormatInt(id, 10)+" user id: "+
				strconv.FormatInt(uid, 10), zap.Error(err2))
		}
	}()
	return nil
}

func (i *interactiveService) Unlike(ctx context.Context, biz string, id int64, uid int64) error {
	err := i.repo.DecreaseLikeCount(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	// 异步发送取消点赞事件
	go func() {
		err2 := i.likeRankEvent.ProduceLikeRankEvent(interactive.LikeRankEvent{
			ArticleID: id,
			Change:    true,
		})
		if err2 != nil {
			i.log.Error("produce unlike rank event failed, article id: "+
				strconv.FormatInt(id, 10)+" user id: "+
				strconv.FormatInt(uid, 10), zap.Error(err2))
		}
	}()
	return nil
}

func NewInteractiveService(repo repository.InteractiveRepository, articleRepository repository.ArticleRepository, event interactive.LikeRankEventProducer, logger *zap.Logger) InteractiveService {
	return &interactiveService{repo: repo, articleRepo: articleRepository, likeRankEvent: event, log: logger}
}

func (i *interactiveService) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncreaseReadCount(ctx, biz, bizId)
}
