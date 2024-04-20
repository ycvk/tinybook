package ioc

import (
	"geek_homework/tinybook/internal/job"
	"geek_homework/tinybook/internal/service"
	"github.com/bsm/redislock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

func InitRankingJob(svc service.RankingService, client *redislock.Client, logger *zap.Logger) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30, client, logger)
}

func InitJobs(log *zap.Logger, rankJob *job.RankingJob) *cron.Cron {
	builder := job.NewCronJobBuilder(log, prometheus.SummaryOpts{
		Namespace: "tinybook",
		Subsystem: "job",
		Name:      "cron_job",
		Help:      "定时任务耗时统计",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	c := cron.New(cron.WithSeconds())
	_, err := c.AddJob("@every 60m", builder.Build(rankJob)) //添加ranking job
	if err != nil {
		panic(err)
	}
	return c
}
