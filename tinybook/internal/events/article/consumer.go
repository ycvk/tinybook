package article

import (
	"context"
	"geek_homework/tinybook/internal/events"
	"geek_homework/tinybook/internal/repository"
	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"time"
)

const GroupArticleRead = "group-article-read"

type KafkaConsumer struct {
	reader *kafka.Reader
	repo   repository.InteractiveRepository
	log    *zap.Logger
}

func NewKafkaConsumer(repo repository.InteractiveRepository, log *zap.Logger) *KafkaConsumer {
	reader := InitReader(GroupArticleRead, TopicArticleRead)
	return &KafkaConsumer{
		repo:   repo,
		log:    log,
		reader: reader,
	}
}

func (k *KafkaConsumer) Start() {
	go func() {
		ctx := context.Background()
		k.Consume(ctx)
	}()
}

func (k *KafkaConsumer) Consume(ctx context.Context) {
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			k.log.Error("close kafka consumer failed", zap.Error(err))
		}
	}(k.reader)
	for {
		// 读取消息
		message, err := k.reader.ReadMessage(ctx)
		if err != nil {
			k.log.Error("read message failed", zap.Error(err))
			continue
		}
		// 解析消息
		var event ReadEvent
		err = sonic.Unmarshal(message.Value, &event)
		if err != nil {
			k.log.Error("consumer unmarshal message failed", zap.Error(err))
			continue
		}
		// 业务逻辑
		err = k.repo.IncreaseReadCount(ctx, "article", event.ArticleID)
		if err != nil {
			k.log.Error("consumer increase read count failed", zap.Error(err))
			continue
		}
	}
}

func CollectConsumer(consumer *KafkaConsumer) []events.Consumer {
	return []events.Consumer{consumer}
}

func InitReader(groupId string, topic string) *kafka.Reader {
	type config struct {
		Brokers string `yaml:"brokers"`
	}
	var cfg config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	split := strings.Split(cfg.Brokers, ",")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        split,
		GroupID:        groupId,
		Topic:          topic,
		MinBytes:       10e3,                   // 10KB
		MaxBytes:       10e6,                   // 10MB
		MaxWait:        500 * time.Millisecond, // 500ms内有数据就返回, 即使没达到MinBytes, 与MinBytes取最小值, 有一个满足就返回
		CommitInterval: time.Second,            // 多久自动commit一次offset
		StartOffset:    kafka.LastOffset,       // 从最新的offset开始读取
	})
	return r
}
