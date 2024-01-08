package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Job struct {
	Id       int64
	Status   int
	Version  int
	Name     string `gorm:"column:name;type:varchar(255);not null;uniqueIndex"`
	Executor string

	NextRunTime int64 `gorm:"column:next_run_time;index"` // 下次执行时间
	Ctime       int64
	Utime       int64
}

const (
	JobStatusUnknown = iota
	JobStatusWaiting
	JobStatusRunning
	JobStatusDone
)

type CronJobDao interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}

type GormCronJobDao struct {
	db *gorm.DB
}

func NewGormCronJobDao(db *gorm.DB) CronJobDao {
	return &GormCronJobDao{db: db}
}

func (g *GormCronJobDao) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	now := time.Now()
	db := g.db.WithContext(ctx)
	return db.Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"next_run_time": t.UnixMilli(),
		"utime":         now.UnixMilli(),
	}).Error
}

func (g *GormCronJobDao) Preempt(ctx context.Context) (Job, error) {
	db := g.db.WithContext(ctx)
	for {
		var job Job
		now := time.Now().UnixMilli()
		err := db.Where("status = ? and next_run_time < ?", JobStatusWaiting, now).First(&job).Error
		if err != nil {
			return Job{}, err
		}
		err = db.Model(&job).Where("version = ?", job.Version).
			Updates(map[string]any{
				"status":  JobStatusRunning,
				"version": job.Version + 1,
				"utime":   now,
			}).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { // 乐观锁失败，说明已经被其他协程抢占
				continue
			}
			return Job{}, err
		}
		return job, nil
	}
}

func (g *GormCronJobDao) Release(ctx context.Context, id int64) error {
	db := g.db.WithContext(ctx)
	now := time.Now().UnixMilli()
	err := db.Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"status": JobStatusWaiting,
		"utime":  now,
	}).Error
	return err
}

func (g *GormCronJobDao) UpdateUTime(ctx context.Context, id int64) error {
	db := g.db.WithContext(ctx)
	now := time.Now().UnixMilli()
	err := db.Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"utime": now,
	}).Error
	return err
}
