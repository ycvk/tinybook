package kafkax

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/kafka-go"
	"sync"
)

// ReaderCollector 是一个自定义的结构体，用于收集 Kafka reader的统计信息
type ReaderCollector struct {
	Reader    *kafka.Reader // kafka.Reader 是 kafka-go 库中用于读取 Kafka 消息的结构体
	NameSpace string
	Subsystem string
	mu        sync.Mutex // 互斥锁，用于确保并发安全

}

func NewReaderCollector(reader *kafka.Reader) *ReaderCollector {
	return &ReaderCollector{
		Reader:    reader,
		NameSpace: "kafka",
		Subsystem: "reader",
	}
}

// 定义 Prometheus 监控指标的标签
var readerLabels = []string{"client_id", "group_id", "topic", "partition"}

// 下面是多个 Prometheus 监控指标的定义
// 使用 CounterVec 和 GaugeVec 来定义计数器和仪表盘类型的指标
// 计数器用于统计发生次数，仪表盘用于显示当前值
var (
	// 以下是计数器类指标，使用 CounterVec 来记录计数器

	// Kafka reader的连接尝试次数
	dialsDescription = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_reader_dial", Help: "kafka消费者的连接次数"}, readerLabels)

	// Kafka reader的拉取操作次数
	fetchesDescription = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_reader_fetch", Help: "kafka消费者的拉取次数"}, readerLabels)

	// Kafka reader读取的消息数量
	messagesDescription = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_reader_message", Help: "kafka消费者读取消息的数量"}, readerLabels)

	// Kafka reader读取的字节总量
	bytesDescription = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_reader_message_bytes", Help: "kafka消费者读取的字节总量"}, readerLabels)

	// Kafka reader的重平衡次数
	rebalancedDescription = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_reader_rebalanced", Help: "kafka消费者重balance次数"}, readerLabels)

	// Kafka reader的超时次数
	timeoutsDescription = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_reader_timeout", Help: "kafka消费者超时次数"}, readerLabels)

	// Kafka reader的错误次数
	errorsDescription = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "kafka_reader_error", Help: "kafka消费者error次数"}, readerLabels)

	// 以下是平均值类指标，使用 GaugeVec 来记录平均值

	// Kafka reader连接时间的平均值
	dialTimeDescriptionAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_dial_seconds_avg", Help: "kafka消费者连接时间平均值"}, readerLabels)

	// Kafka reader读取操作时间的平均值
	readTimeDescriptionAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_read_seconds_avg", Help: "kafka消费者读取时间的平均值"}, readerLabels)

	// Kafka reader等待时间的平均值
	waitTimeDescriptionAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_wait_seconds_avg", Help: "kafka消费者等待时间的平均值"}, readerLabels)

	// Kafka reader每次拉取操作的平均大小
	fetchSizeDescriptionAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_fetch_size_avg", Help: "kafka消费者每次拉取消息的平均大小"}, readerLabels)

	// Kafka reader每次拉取操作的平均字节数
	fetchBytesDescriptionAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_fetch_bytes_avg", Help: "kafka消费者每次拉取消息的平均字节"}, readerLabels)

	// 以下是当前值类指标，使用 GaugeVec 来记录当前值

	// Kafka reader的当前偏移量
	offsetDescription = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_offset", Help: "kafka消费者当前offset"}, readerLabels)

	// Kafka reader的当前延迟
	lagDescription = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_lag", Help: "kafka消费者当前延迟"}, readerLabels)

	// Kafka reader配置的最小拉取字节
	//minBytesDescription = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_config_fetch_bytes_min", Help: "kafka消费者配置的拉取的min_bytes"}, readerLabels)

	// Kafka reader配置的最大拉取字节
	//maxBytesDescription = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_config_fetch_bytes_max", Help: "kafka消费者配置的拉取的max_bytes"}, readerLabels)

	// Kafka reader的最大等待时间
	maxWaitDescription = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_fetch_wait_max", Help: "kafka消费者max_wait"}, readerLabels)

	// Kafka reader内部内存队列的长度
	queueLengthDescription = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_queue_length", Help: "内部内存队列的长度"}, readerLabels)

	// Kafka reader内部内存队列的容量
	queueCapacityDescription = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "kafka_reader_queue_capacity", Help: "内部内存队列的容量"}, readerLabels)
)

// Describe 方法留空，因为所有指标都是用 New...Vec 创建的，Prometheus 会自动处理它们的描述
func (r *ReaderCollector) Describe(_ chan<- *prometheus.Desc) {
}

// readerCounter 函数用于更新计数器类型的指标
func readerCounter(counter *prometheus.CounterVec, value float64, id string, groupId string, topic string, partition string) prometheus.Counter {
	counterWithLabel := counter.WithLabelValues(id, groupId, topic, partition)
	counterWithLabel.Add(value)
	return counterWithLabel
}

// readerGauge 函数用于更新仪表盘类型的指标
func readerGauge(gauge *prometheus.GaugeVec, value float64, id string, groupId string, topic string, partition string) prometheus.Metric {
	gaugeWithLabel := gauge.WithLabelValues(id, groupId, topic, partition)
	if value != 0 {
		gaugeWithLabel.Set(value)
	}
	return gaugeWithLabel
}

// Collect 方法用于收集指标数据并发送到 Prometheus
func (r *ReaderCollector) Collect(metrics chan<- prometheus.Metric) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stats := r.Reader.Stats() // 获取 Kafka reader的统计信息
	groupID := r.Reader.Config().GroupID
	clientId, topic, partition := stats.ClientID, stats.Topic, stats.Partition

	// 使用 readerCounter 和 readerGauge 函数更新指标，并将它们发送到 Prometheus
	metrics <- readerCounter(dialsDescription, float64(stats.Dials), clientId, groupID, topic, partition)
	metrics <- readerCounter(fetchesDescription, float64(stats.Fetches), clientId, groupID, topic, partition)
	metrics <- readerCounter(messagesDescription, float64(stats.Messages), clientId, groupID, topic, partition)
	metrics <- readerCounter(bytesDescription, float64(stats.Bytes), clientId, groupID, topic, partition)
	metrics <- readerCounter(rebalancedDescription, float64(stats.Rebalances), clientId, groupID, topic, partition)
	metrics <- readerCounter(timeoutsDescription, float64(stats.Timeouts), clientId, groupID, topic, partition)
	metrics <- readerCounter(errorsDescription, float64(stats.Errors), clientId, groupID, topic, partition)
	metrics <- readerGauge(dialTimeDescriptionAvg, stats.DialTime.Avg.Seconds(), clientId, groupID, topic, partition)
	metrics <- readerGauge(readTimeDescriptionAvg, stats.ReadTime.Avg.Seconds(), clientId, groupID, topic, partition)
	metrics <- readerGauge(waitTimeDescriptionAvg, stats.WaitTime.Avg.Seconds(), clientId, groupID, topic, partition)
	metrics <- readerGauge(fetchSizeDescriptionAvg, float64(stats.FetchSize.Avg), clientId, groupID, topic, partition)
	metrics <- readerGauge(fetchBytesDescriptionAvg, float64(stats.FetchBytes.Avg), clientId, groupID, topic, partition)
	metrics <- readerGauge(offsetDescription, float64(stats.Offset), clientId, groupID, topic, partition)
	metrics <- readerGauge(lagDescription, float64(stats.Lag), clientId, groupID, topic, partition)
	//metrics <- readerGauge(minBytesDescription, float64(stats.MinBytes), clientId, groupID, topic, partition)
	//metrics <- readerGauge(maxBytesDescription, float64(stats.MaxBytes), clientId, groupID, topic, partition)
	metrics <- readerGauge(maxWaitDescription, float64(stats.MaxWait), clientId, groupID, topic, partition)
	metrics <- readerGauge(queueLengthDescription, float64(stats.QueueLength), clientId, groupID, topic, partition)
	metrics <- readerGauge(queueCapacityDescription, float64(stats.QueueCapacity), clientId, groupID, topic, partition)
}
