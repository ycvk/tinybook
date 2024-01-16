package service

import (
	"context"
	"github.com/samber/lo"
	"math"
	"time"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/internal/domain"
	"tinybook/tinybook/internal/repository"
	"tinybook/tinybook/pkg/priorityqueue"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	InteractiveSvc intrv1.InteractiveServiceClient
	ArticleSvc     ArticleService
	BatchSize      int // 每次获取的文章数量
	topNum         int // 排行榜数量
	ScoreFunc      func(likeCount int64, utime time.Time) float64
	queue          *priorityqueue.PriorityQueue[domain.Article, float64] // 优先队列
	rankingRepo    repository.RankingRepository
}

func (b *BatchRankingService) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return b.rankingRepo.GetTopN(ctx)
}

func NewBatchRankingService(interactiveSvc intrv1.InteractiveServiceClient, articleSvc ArticleService, repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		InteractiveSvc: interactiveSvc,
		ArticleSvc:     articleSvc,
		BatchSize:      1000,
		topNum:         100,
		rankingRepo:    repo,
		ScoreFunc: func(likeCount int64, utime time.Time) float64 {
			return float64(likeCount-1) / math.Pow(time.Now().Sub(utime).Seconds()+2, 1.8)
		},
		queue: priorityqueue.New[domain.Article, float64](priorityqueue.MinHeap),
	}
}

func (b *BatchRankingService) TopN(ctx context.Context) error {
	topN, err := b.topN(ctx)
	if err != nil {
		return err
	}
	// 放到缓存中
	return b.rankingRepo.ReplaceTopN(ctx, topN)
}

func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	now := time.Now()
	ddl := now.Add(-time.Hour * 24 * 7) // 一周前
	offset := 0
	for {
		// 获取article
		listPub, err := b.ArticleSvc.ListPub(ctx, now, b.BatchSize, offset)
		if err != nil {
			return nil, err
		}
		if len(listPub) == 0 {
			break
		}
		// 获取article的id
		ids := lo.Map(listPub, func(item domain.Article, index int) int64 {
			return item.ID
		})
		// 根据id获取article的interactive
		byIds, err := b.InteractiveSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
			Biz: "article",
			Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		interactives := byIds.GetInteractives()
		for _, article := range listPub {
			intr := interactives[article.ID]
			score := b.ScoreFunc(intr.LikeCount, time.Unix(article.Utime, 0))
			if b.queue.Len() > b.topNum {
				// 队列已满
				// 比较队列中最小的元素和当前元素的大小
				// 如果当前元素大于最小元素，则替换最小元素
				// 否则，跳过当前元素
				queueMin := b.queue.Get()
				if score < queueMin.Priority { // 当前元素小于最小元素 放回去
					b.queue.PutItem(queueMin)
					continue
				}
				b.queue.Put(article, score)
			}
			b.queue.Put(article, score)
		}
		offset += len(listPub)
		if len(listPub) < b.BatchSize || listPub[len(listPub)-1].Utime < ddl.Unix() { // 如果最后一条数据的时间超过了ddl，就不再继续获取
			break
		}
	}
	articles := make([]domain.Article, b.queue.Len())
	for b.queue.Len() > 0 {
		articles = append(articles, b.queue.Get().Value)
	}
	return articles, nil
}
