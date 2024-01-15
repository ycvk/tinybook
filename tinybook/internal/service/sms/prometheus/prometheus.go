package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"tinybook/tinybook/internal/service/sms"
)

type Decorator struct {
	service sms.Service
	vec     *prometheus.SummaryVec
}

// NewDecorator 创建一个装饰器 用于记录耗时
func NewDecorator(service sms.Service) *Decorator {
	return &Decorator{
		service: service,
		vec: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "tinybook",
			Subsystem: "sms",
			Name:      "send_code_duration",
			Help:      "发送验证码耗时",
		}, []string{"status"}),
	}
}

func (d *Decorator) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 记录耗时
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		d.vec.WithLabelValues(tplId).Observe(float64(duration))
	}()
	return d.service.Send(ctx, tplId, args, numbers...)
}
