package repository

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository/dao"
	"time"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, jId int64) error
	UpdateUTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, time time.Time) error
}

type cronJobRepository struct {
	dao dao.CronJobDao
}

func NewCronJobRepository(dao dao.CronJobDao) CronJobRepository {
	return &cronJobRepository{dao: dao}

}

func (c *cronJobRepository) UpdateNextTime(ctx context.Context, id int64, time time.Time) error {
	return c.dao.UpdateNextTime(ctx, id, time)
}

func (c *cronJobRepository) UpdateUTime(ctx context.Context, id int64) error {
	return c.dao.UpdateUTime(ctx, id)
}

func (c *cronJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := c.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		Id: job.Id,
	}, nil
}

func (c *cronJobRepository) Release(ctx context.Context, jId int64) error {
	return c.dao.Release(ctx, jId)
}
