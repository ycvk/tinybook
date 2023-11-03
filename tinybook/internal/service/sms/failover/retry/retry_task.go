package retry

import (
	"github.com/cockroachdb/errors"
	"log/slog"
	"math"
	"math/rand"
	"time"
)

type AsyncRetry interface {
	StartRetryLoop(task Task) (bool, error)
	RecordResult(success bool)
	CheckErrorRate() bool
}

type RetryTask struct {
	MaxRetries          int           // 最大重试次数
	BaseInterval        time.Duration // 基础间隔时间
	Multiplier          float64       // 间隔增加的倍数
	RandomizationFactor float64       // 随机化因素，用于避免网络拥塞
	ErrRateMonitor      ErrorMonitor  // 错误率监控器
}

func NewRetryTask(max int, monitor ErrorMonitor) AsyncRetry {
	return &RetryTask{
		MaxRetries:          max,                    // 最大重试次数
		BaseInterval:        500 * time.Millisecond, // 基础间隔时间为500毫秒
		Multiplier:          2,                      // 间隔增加的倍数，以2为例，在无随机情况下，每次重试的间隔时间为：0.5秒、1秒、2秒、4秒
		RandomizationFactor: 0.5,                    // 随机化因素，用于避免网络拥塞，也为防止雷鸣效应，0.5表示实际间隔将随机地增加或减少最多50%
		ErrRateMonitor:      monitor,
	}
}

// Task 代表一个可以重试的工作单元
type Task func() error

// exponentialBackoff 采用指数退避 + 随机因子策略，根据重试次数计算下一次重试前的等待时间
func exponentialBackoff(retries int, config *RetryTask) time.Duration {
	// 计算本次重试的间隔时间 = 基础间隔时间 * 间隔增加的倍数^重试次数
	expInterval := float64(config.BaseInterval) * math.Pow(config.Multiplier, float64(retries))
	// 应用随机因子 = 间隔时间 * (1 - 随机化因子 + 随机化因子 * 2 * 随机数)
	expInterval = expInterval * (1 - config.RandomizationFactor + (rand.Float64() * 2 * config.RandomizationFactor))
	// 限制最大间隔时间，确保间隔在规定范围内
	interval := time.Duration(expInterval)
	interval = max(interval, config.BaseInterval)
	return interval
}

// StartRetryLoop 启动一个循环,根据重试策略重试任务
func (r *RetryTask) StartRetryLoop(task Task) (bool, error) {
	retries := 0
	for {
		err := task()
		if err != nil {
			// 任务执行失败，记录结果
			r.ErrRateMonitor.RecordResult(false)
			if retries < r.MaxRetries {
				// 如果任务执行失败，且重试次数未达到最大重试次数，等待一段时间后重试
				time.Sleep(exponentialBackoff(retries, r))
				retries++
			} else {
				// 重试次数达到最大重试次数，任务执行失败
				slog.Error("最大重试后任务失败", "maxRetries", r.MaxRetries)
				return false, errors.New("最大重试后任务失败")
			}
		} else {
			// 任务执行成功，记录结果
			r.ErrRateMonitor.RecordResult(true)
			// 任务执行成功，退出循环
			return true, nil
		}
	}
}

func (r *RetryTask) RecordResult(success bool) {
	// 记录结果
	r.ErrRateMonitor.RecordResult(success)
}

func (r *RetryTask) CheckErrorRate() bool {
	// 检查错误率是否超过阈值
	return r.ErrRateMonitor.CheckErrorRate()
}
