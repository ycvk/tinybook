package job

import (
	"context"
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"time"
	"tinybook/tinybook/internal/domain"
	"tinybook/tinybook/internal/service"
)

type Executor interface {
	Name() string
	Execute(ctx context.Context, job domain.Job) error
}

// LocalFuncExecutor 本地函数执行器
type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, job domain.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: make(map[string]func(ctx context.Context, job domain.Job) error),
	}
}

func (l *LocalFuncExecutor) Name() string {
	return "local_func"
}

func (l *LocalFuncExecutor) Execute(ctx context.Context, job domain.Job) error {
	f, ok := l.funcs[job.Name]
	if !ok {
		return errors.New("unknown job name: " + job.Name)
	}
	return f(ctx, job)
}

func (l *LocalFuncExecutor) RegisterFunc(name string, f func(ctx context.Context, job domain.Job) error) {
	if l.funcs == nil {
		l.funcs = make(map[string]func(ctx context.Context, job domain.Job) error)
	}
	l.funcs[name] = f

}

type Scheduler struct {
	dbTimeout time.Duration
	service   service.CronJobService
	executors map[string]Executor
	log       *zap.Logger
	limiter   *semaphore.Weighted
}

func NewScheduler(jobService service.CronJobService, log *zap.Logger) *Scheduler {
	return &Scheduler{
		service:   jobService,
		log:       log,
		dbTimeout: time.Second,
		limiter:   semaphore.NewWeighted(100), // 限制并发数
		executors: make(map[string]Executor),
	}
}

func (s *Scheduler) RegisterExecutor(executor Executor) {
	s.executors[executor.Name()] = executor
}

func (s *Scheduler) Schedule(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}
		err2 := s.limiter.Acquire(ctx, 1) // 从信号量中获取一个信号量
		if err2 != nil {
			s.log.Error("acquire semaphore failed", zap.Error(err2))
			return
		}
		timeout, cancelFunc := context.WithTimeout(ctx, s.dbTimeout)
		// 从数据库中获取任务
		job, err := s.service.Preempt(timeout)
		cancelFunc()
		if err != nil {
			continue
		}
		// 找到对应的 executor
		executor, ok := s.executors[job.Executor]
		if !ok { // 未知的 executor，直接跳过
			s.log.Error("unknown executor", zap.String("executor", job.Executor), zap.Int64("job_id", job.Id))
			continue
		}
		go func() {
			// 任务执行完毕，释放资源
			defer func() {
				s.limiter.Release(1) // 释放一个信号量
				job.CancelFunc()
			}()
			// 执行任务
			err2 := executor.Execute(ctx, job)
			if err != nil {
				s.log.Error("execute job failed", zap.Error(err2), zap.Int64("job_id", job.Id))
			}
			return
		}()
	}
}
