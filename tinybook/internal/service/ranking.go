package service

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/pkg/priorityqueue"
	"github.com/samber/lo"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
}

type BatchRankingService struct {
	InteractiveSvc InteractiveService
	ArticleSvc     ArticleService
	BatchSize      int // 每次获取的文章数量
	topNum         int // 排行榜数量
	ScoreFunc      func(likeCount int64, utime time.Time) float64
	queue          *priorityqueue.PriorityQueue[domain.Article, float64]
}

func NewBatchRankingService(interactiveSvc InteractiveService, articleSvc ArticleService) RankingService {
	return &BatchRankingService{
		InteractiveSvc: interactiveSvc,
		ArticleSvc:     articleSvc,
		BatchSize:      1000,
		topNum:         100,
		ScoreFunc: func(likeCount int64, utime time.Time) float64 {
			return float64(likeCount-1) / math.Pow(time.Now().Sub(utime).Seconds()+2, 1.8)
		},
		queue: priorityqueue.New[domain.Article, float64](priorityqueue.MinHeap),
	}
}

func (b *BatchRankingService) TopN(ctx context.Context) error {
	panic("implement me")
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
		byIds, err := b.InteractiveSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		for _, article := range listPub {
			intr := byIds[article.ID]
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
