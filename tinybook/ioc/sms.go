package ioc

import (
	"github.com/redis/go-redis/v9"
	"time"
	"tinybook/tinybook/internal/repository"
	"tinybook/tinybook/internal/service/sms"
	"tinybook/tinybook/internal/service/sms/failover/retry"
	"tinybook/tinybook/internal/service/sms/localsms"
	"tinybook/tinybook/pkg/limiter"
)

func InitSMSService(cmd redis.Cmdable, repo repository.SMSRepository) sms.Service {
	// 错误率监控器, 当 30 秒内错误率超过 30% 时, 将触发重试任务
	monitor := retry.NewErrorRateMonitor(0.3, 0.5, 30*time.Second)

	return retry.NewAsyncFailoverSMSService(
		// 限流器, 30 秒钟最多发送 10 条短信
		limiter.NewRedisSlideWindowLimiter(cmd, 30*time.Second, 10),
		// 本地短信服务
		localsms.NewService(),
		// 短信存储库
		repo,
		// 重试任务, 最大重试次数3次
		retry.NewRetryTask(3, monitor),
	)
	//return localsms.NewService()
}
