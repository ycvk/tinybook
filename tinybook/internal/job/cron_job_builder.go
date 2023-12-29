package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	log    *zap.Logger
	vector *prometheus.SummaryVec
}

func NewCronJobBuilder(log *zap.Logger, opts prometheus.SummaryOpts) *CronJobBuilder {
	vector := prometheus.NewSummaryVec(opts, []string{"job", "success"})
	return &CronJobBuilder{log: log, vector: vector}
}

func (c *CronJobBuilder) Build(job Job) cron.Job {
	return cron.FuncJob(func() {
		name := job.Name()
		start := time.Now()
		err := job.Run()
		if err != nil {
			c.log.Error("job run failed", zap.String("job", name), zap.Error(err))
		}
		end := time.Since(start).Milliseconds()
		c.vector.WithLabelValues(name, strconv.FormatBool(err == nil)).Observe(float64(end))
	})
}
