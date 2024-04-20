package kafkax

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/kafka-go"
	"sync"
)

type WriterCollector struct {
	Writer    *kafka.Writer
	Namespace string
	Subsystem string
	mu        sync.Mutex // 互斥锁，用于确保并发安全
}

func NewWriterCollector(writer *kafka.Writer) *WriterCollector {
	return &WriterCollector{Writer: writer, Namespace: "kafka", Subsystem: "writer"}
}

// 定义 Prometheus 监控指标的标签
var writerLabels = []string{"client_id", "topic"}

var (
	writerMessageCount = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_writer_message_count", Help: "kafka生产者生产消息数量"}, writerLabels)
	writerBatchBytes   = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_writer_batch_bytes", Help: "kafka生产者发送的字节数"}, writerLabels)
	writerBatchTime    = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_writer_batch_time", Help: "kafka生产者发送的时间"}, writerLabels)
)

// readerCounter 函数用于更新计数器类型的指标
func writerCounter(counter *prometheus.CounterVec, value float64, id string, topic string) prometheus.Counter {
	counterWithLabel := counter.WithLabelValues(id, topic)
	counterWithLabel.Add(value)
	return counterWithLabel
}

// readerGauge 函数用于更新仪表盘类型的指标
func writerGauge(gauge *prometheus.GaugeVec, value float64, id string, topic string) prometheus.Metric {
	gaugeWithLabel := gauge.WithLabelValues(id, topic)
	if value != 0 {
		gaugeWithLabel.Set(value)
	}
	return gaugeWithLabel
}

func (w *WriterCollector) Describe(_ chan<- *prometheus.Desc) {
}

func (w *WriterCollector) Collect(metrics chan<- prometheus.Metric) {
	w.mu.Lock()
	defer w.mu.Unlock()
	stats := w.Writer.Stats()
	clientId, topic := stats.ClientID, stats.Topic

	metrics <- writerCounter(writerMessageCount, float64(stats.Messages), clientId, topic)
	metrics <- writerGauge(writerBatchTime, stats.BatchTime.Sum.Seconds(), clientId, topic)
	metrics <- writerGauge(writerBatchBytes, float64(stats.BatchBytes.Sum), clientId, topic)
}
