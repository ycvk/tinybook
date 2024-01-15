package ratelimit

import (
	"context"
	"github.com/pingcap/errors"
	"tinybook/tinybook/internal/service/sms"
	"tinybook/tinybook/pkg/limiter"
)

var ErrLimitReached = errors.New("rate limit reached in sms service")

// RateLimitSMSService 限流短信服务 使用组合模式
//
//	type RateLimitSMSService struct {
//		limiter limiter.Limiter
//		sms.Service
//		key string
//	}
//
// RateLimitSMSService 限流短信服务 使用装饰器模式
// 装饰器模式比起组合模式，可以阻止绕开限流器直接调用短信服务，但必须实现短信服务的所有接口
// 组合模式可以只实现部分接口，但是可以直接调用短信服务，绕开限流器
type RateLimitSMSService struct {
	limiter limiter.Limiter
	service sms.Service
	key     string
}

func NewRateLimitSMSService(limiter limiter.Limiter, service sms.Service) *RateLimitSMSService {
	return &RateLimitSMSService{
		limiter: limiter,
		service: service,
		key:     "sms-limiter",
	}
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limit, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limit {
		return ErrLimitReached
	}
	return r.service.Send(ctx, tplId, args, numbers...)
}
