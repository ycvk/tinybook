package service

import (
	"context"
	"go.uber.org/zap"
	"time"
	"tinybook/tinybook/internal/domain"
	"tinybook/tinybook/internal/repository"
)

type CronJobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextRunTime(ctx context.Context, j domain.Job) error
}

type cronJobService struct {
	log             *zap.Logger
	repo            repository.CronJobRepository
	refreshInterval time.Duration
}

func NewCronJobService(log *zap.Logger, repo repository.CronJobRepository) CronJobService {
	return &cronJobService{
		log:             log,
		repo:            repo,
		refreshInterval: time.Second * 60, // 每隔一分钟刷新一次 job 的更新时间
	}

}

func (c *cronJobService) ResetNextRunTime(ctx context.Context, j domain.Job) error {
	nextTime := j.NextTime()
	return c.repo.UpdateNextTime(ctx, j.Id, nextTime)
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.refresh(job.Id) // 每隔一段时间刷新一次 job 的更新时间
			}
		}
	}()

	job.CancelFunc = func() {
		ticker.Stop()
		withTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		err2 := c.repo.Release(withTimeout, job.Id)
		if err2 != nil {
			c.log.Error("release job failed", zap.Error(err2))
		}
	}
	return job, err
}

func (c *cronJobService) refresh(id int64) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	err := c.repo.UpdateUTime(timeout, id)
	if err != nil {
		c.log.Error("update job utime failed", zap.Error(err))
	}
}
