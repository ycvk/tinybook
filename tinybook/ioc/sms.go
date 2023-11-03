package ioc

import (
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/service/sms"
	"geek_homework/tinybook/internal/service/sms/failover/retry"
	"geek_homework/tinybook/internal/service/sms/localsms"
	"geek_homework/tinybook/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSMSService(cmd redis.Cmdable, repo repository.SMSRepository) sms.Service {
	return retry.NewAsyncFailoverSMSService(
		// 限流器, 30 秒钟最多发送 10 条短信
		limiter.NewRedisSlideWindowLimiter(cmd, 30*time.Second, 10),
		// 本地短信服务
		localsms.NewService(),
		// 短信存储库
		repo,
		// 错误率监控器, 当 30 秒内错误率超过 30% 时, 将触发重试任务
		retry.NewErrorRateMonitor(0.3, 0.5, 30*time.Second),
		// 重试任务, 最大重试次数3次
		retry.NewRetryTask(),
	)
	//return localsms.NewService()
}
