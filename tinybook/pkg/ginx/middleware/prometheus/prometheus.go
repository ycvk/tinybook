package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Builder struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceId string
}

func NewBuilder(namespace, subsystem, name, help, instanceId string) *Builder {
	return &Builder{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       name,
		Help:       help,
		InstanceId: instanceId,
	}
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	labels := []string{"method", "path", "status"} // 用于标识请求的标签 例如：method=GET path=/api/v1/user status=200
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_response_time",
		Help:      b.Help,
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)
	prometheus.MustRegister(vector) // 注册到 prometheus
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			path := ctx.FullPath()
			status := ctx.Writer.Status()
			// 根据标签获取对应的 Summary 对象，然后将请求耗时记录到 Summary 中
			vector.WithLabelValues(method, path, strconv.Itoa(status)).Observe(float64(duration))
		}()
		ctx.Next()
	}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_active_request",
		Help:      b.Help,
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()       // 每次请求进来，就将 Gauge 加 1
		defer gauge.Dec() // 每次请求结束，就将 Gauge 减 1
		ctx.Next()
	}
}
