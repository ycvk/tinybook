package retry

import (
	"context"
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/service/sms"
	"geek_homework/tinybook/pkg/limiter"
	"github.com/cockroachdb/errors"
	"log/slog"
)

type AsyncFailoverSMSService struct {
	services sms.Service
	smsRepo  repository.SMSRepository
	limiter  limiter.Limiter // 限流器
	limitKey string          // 限流器的key
	retryTh  AsyncRetry      // 重试任务
}

func NewAsyncFailoverSMSService(limiter limiter.Limiter, services sms.Service, smsRe repository.SMSRepository, task AsyncRetry) sms.Service {
	return &AsyncFailoverSMSService{
		services: services,
		smsRepo:  smsRe,
		retryTh:  task,
		limiter:  limiter,
		limitKey: "failover_async_sms", // 限流器的key
	}
}

func (f *AsyncFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 定义一个重试函数，用于重试发送短信
	retryFunc := func() error {
		return f.services.Send(ctx, tplId, args, numbers...)
	}

	limit, err := f.limiter.Limit(ctx, f.limitKey)
	if err != nil {
		return err
	}
	// 如果限流，将请求转储到数据库，后续再另外启动一个 goroutine 异步发送出去
	if limit {
		err := f.smsRepo.Save(ctx, tplId, args, numbers...)
		if err != nil {
			// 存储到数据库失败，直接返回error
			slog.Error("存储发送请求到数据库失败", "err", err)
			return err
		}
		// 启动一个 goroutine 去异步调用重试函数
		go func() {
			ok, retryErr := f.retryTh.StartRetryLoop(retryFunc)
			if !ok || retryErr != nil {
				slog.Error("重试发送短信失败", "err", retryErr)
			} else {
				slog.Info("重试发送短信成功")
				// 重试成功后，删除数据库中的记录
				delErr := f.smsRepo.Delete(ctx, tplId, args, numbers...)
				if delErr != nil {
					slog.Error("删除数据库中的记录失败", "err", delErr)
				}
			}
		}()
		return err
	}
	// 如果没有限流，直接发送等待结果
	err = f.services.Send(ctx, tplId, args, numbers...)
	// 记录发送结果到错误率监控器
	f.retryTh.RecordResult(errors.Is(err, nil))
	// 检查错误率是否超过阈值
	rate := f.retryTh.CheckErrorRate()
	if rate {
		// 如果超过阈值，将请求转储到数据库，后续再另外启动一个 goroutine 异步发送出去
		err := f.smsRepo.Save(ctx, tplId, args, numbers...)
		if err != nil {
			// 存储到数据库失败，直接返回error
			slog.Error("存储发送请求到数据库失败", "err", err)
			return err
		}
		// 启动一个 goroutine 去异步调用重试函数
		go func() {
			ok, retryErr := f.retryTh.StartRetryLoop(retryFunc)
			if !ok || retryErr != nil {
				slog.Error("重试发送短信失败", "err", retryErr)
			} else {
				slog.Info("重试发送短信成功")
				// 重试成功后，删除数据库中的记录
				delErr := f.smsRepo.Delete(ctx, tplId, args, numbers...)
				if delErr != nil {
					slog.Error("删除数据库中的记录失败", "err", delErr)
				}
			}
		}()
		return nil
	}
	return err
}
