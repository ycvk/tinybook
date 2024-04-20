package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"tinybook/tinybook/internal/domain"
	"tinybook/tinybook/internal/service/oauth2/wechat"
)

type Decorator struct {
	wechat.Service
	sum prometheus.Summary
}

// NewDecorator 组合模式 用于记录耗时
func NewDecorator(service wechat.Service, sum prometheus.Summary) *Decorator {
	return &Decorator{
		Service: service,
		sum:     sum,
	}
}

func (d *Decorator) Verify(ctx context.Context, code string) (domain.WechatInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		d.sum.Observe(float64(duration))
	}()
	return d.Service.Verify(ctx, code)
}
