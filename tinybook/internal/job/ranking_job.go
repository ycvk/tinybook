package job

import (
	"context"
	"geek_homework/tinybook/internal/service"
	"time"
)

type RankingJob struct {
	rankingSvc service.RankingService
	time       time.Duration
}

func NewRankingJob(rankingSvc service.RankingService, time time.Duration) *RankingJob {
	return &RankingJob{rankingSvc: rankingSvc, time: time}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), r.time)
	defer cancelFunc()
	return r.rankingSvc.TopN(ctx)
}
