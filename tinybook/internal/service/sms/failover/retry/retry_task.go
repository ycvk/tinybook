package retry

import (
	"github.com/cockroachdb/errors"
	"log/slog"
	"math"
	"math/rand"
	"time"
)

type RetryTask struct {
	MaxRetries          int           // 最大重试次数
	BaseInterval        time.Duration // 基础间隔时间
	Multiplier          float64       // 间隔增加的倍数
	RandomizationFactor float64       // 随机化因素，用于避免网络拥塞
}

func NewRetryTask() *RetryTask {
	return &RetryTask{
		MaxRetries:          3,                      // 最大重试次数
		BaseInterval:        500 * time.Millisecond, // 基础间隔时间为500毫秒
		Multiplier:          2,                      // 间隔增加的倍数，以2为例，在无随机情况下，每次重试的间隔时间为：0.5秒、1秒、2秒、4秒
		RandomizationFactor: 0.5,                    // 随机化因素，用于避免网络拥塞，也为防止雷鸣效应，0.5表示实际间隔将随机地增加或减少最多50%
	}
}

// Task represents a unit of work that can be retried
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
func (r *RetryTask) StartRetryLoop(task Task) error {
	retries := 0
	for {
		err := task()
		if err != nil && retries < r.MaxRetries {
			// 如果任务执行失败，且重试次数未达到最大重试次数，等待一段时间后重试
			time.Sleep(exponentialBackoff(retries, r))
			retries++
		} else {
			break
		}
	}
	if retries >= r.MaxRetries {
		slog.Error("最大重试后任务失败", "maxRetries", r.MaxRetries)
		return errors.New("最大重试后任务失败")
	} else {
		// 任务执行成功
		return nil
	}
}
