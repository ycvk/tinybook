package ioc

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitWriter() *kafka.Writer {
	type config struct {
		Brokers string `yaml:"brokers"`
	}
	var cfg config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	split := strings.Split(cfg.Brokers, ",")
	w := &kafka.Writer{
		Addr:                   kafka.TCP(split...),
		BatchTimeout:           1000 * time.Millisecond, // 1秒flush一次
		BatchSize:              10,                      // 10条数据就flush, 与BatchTimeout取最小值, 有一个满足就flush
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,      // 允许自动创建topic
		Compression:            kafka.Lz4, // 压缩算法
	}
	return w
}
