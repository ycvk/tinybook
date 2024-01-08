package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id   int64
	Name string
	// cron 表达式
	Expression string
	Executor   string
	CancelFunc func()
}

// NextTime 计算下次执行时间
func (j *Job) NextTime() time.Time {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, _ := parser.Parse(j.Expression)
	return schedule.Next(time.Now())
}
