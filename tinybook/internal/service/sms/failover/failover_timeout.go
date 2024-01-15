package failover

import (
	"context"
	"errors"
	"sync/atomic"
	"tinybook/tinybook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	services  []sms.Service
	idx       int32 //当前使用的服务的下标
	cnt       int32 //当前服务的失败次数
	threshold int32 //失败次数的阈值
}

func NewTimeoutFailoverSMSService(threshold int32, services ...sms.Service) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{
		services:  services,
		idx:       0,
		cnt:       0,
		threshold: threshold,
	}
}

func (f TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&f.idx)
	cnt := atomic.LoadInt32(&f.cnt)
	// 如果当前服务的失败次数超过了阈值，那么就切换到下一个服务
	if cnt >= f.threshold {
		newIdx := (idx + 1) % int32(len(f.services))
		// 如果切换成功，那么就重置失败次数
		if atomic.CompareAndSwapInt32(&f.idx, idx, newIdx) {
			atomic.StoreInt32(&f.cnt, 0)
		}
	}
	// 获取当前服务的下标
	idx = atomic.LoadInt32(&f.idx)
	// 获取当前服务
	service := f.services[idx]
	// 调用当前服务的发送短信的方法
	err := service.Send(ctx, tplId, args, numbers...)
	switch {
	// 如果没有错误，说明发送成功，直接返回
	case err == nil:
		atomic.StoreInt32(&f.cnt, 0)
		return nil
	// 如果是超时错误，那么就增加失败次数
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt32(&f.cnt, 1)
		return err
	default:
		return err
	}
}
