package job

import (
	"context"
	"geek_homework/tinybook/internal/service"
	"github.com/bsm/redislock"
	"go.uber.org/zap"
	"time"
)

type RankingJob struct {
	log        *zap.Logger
	rankingSvc service.RankingService
	time       time.Duration
	redisLock  *redislock.Client
	key        string
}

func NewRankingJob(rankingSvc service.RankingService, time time.Duration, lock *redislock.Client, l *zap.Logger) *RankingJob {
	return &RankingJob{rankingSvc: rankingSvc, time: time, redisLock: lock, log: l}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	timeout, c := context.WithTimeout(context.Background(), time.Second*3)
	defer c()
	lock, err := r.redisLock.Obtain(timeout, r.key, r.time,
		&redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(time.Millisecond*100), 3), // 重试3次，每次间隔100ms
		})
	if err != nil {
		r.log.Error("ranking job lock failed", zap.Error(err))
		return err
	}
	defer func() {
		withTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
		defer cancelFunc()
		err2 := lock.Release(withTimeout)
		if err2 != nil {
			r.log.Error("ranking job unlock failed", zap.Error(err2))
		}
	}()

	ctx, cancelFunc := context.WithTimeout(context.Background(), r.time)
	defer cancelFunc()
	return r.rankingSvc.TopN(ctx)
}
