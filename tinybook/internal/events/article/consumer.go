package article

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"time"
	"tinybook/tinybook/internal/events"
	"tinybook/tinybook/internal/events/interactive"
	"tinybook/tinybook/internal/repository"
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
		//k.Consume(ctx)
		k.BatchConsume(ctx)
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

func (k *KafkaConsumer) BatchConsume(ctx context.Context) {
	timeWait := k.reader.Config().MaxWait                    // consumer最大等待时间
	msgLen := 10                                             // 一次批量消费的消息数量
	msgs := make([]kafka.Message, 0, msgLen)                 // 一次批量消费的消息 10条
	timeoutCtx, cancel := context.WithTimeout(ctx, timeWait) // 一次批量消费的超时时间
	// 业务逻辑
	fn := func(ms []kafka.Message) error {
		artIds := make([]int64, 0, len(ms))
		for i := range ms {
			if ms[i].Value == nil {
				continue
			}
			// 解析消息
			var event ReadEvent
			er := sonic.Unmarshal(ms[i].Value, &event)
			if er != nil {
				k.log.Error("consumer unmarshal message failed", zap.Error(er))
				continue
			}
			artIds = append(artIds, event.ArticleID)
		}
		// 业务逻辑 批量增加阅读数
		err := k.repo.BatchIncreaseReadCount(ctx, "article", artIds)
		if err != nil {
			k.log.Error("consumer batch increase read count failed", zap.Error(err))
			return err
		}
		return nil
	}

	for {
		msg, err := k.reader.FetchMessage(timeoutCtx) // 从broker获取消息
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) { // 超时
				cancel()
				// 执行业务逻辑 一次批量消费后 提交offset 清空msgs
				err := fn(msgs)
				if err != nil {
					continue
				}
				// 提交offset
				err = k.reader.CommitMessages(ctx, msgs...)
				if err != nil {
					k.log.Error("commit messages failed", zap.Error(err))
					//continue
				}
				msgs = msgs[:0]
				timeoutCtx, cancel = context.WithTimeout(ctx, timeWait) // 重新设置超时时间
				continue
			}
		}
		if len(msgs) == msgLen { // 需要批量消费的消息数量达到10条
			// 执行业务逻辑 一次批量消费后 提交offset 清空msgs
			err := fn(msgs)
			if err != nil {
				continue
			}
			err = k.reader.CommitMessages(ctx, msgs...)
			if err != nil {
				k.log.Error("commit messages failed", zap.Error(err))
			}
			msgs = msgs[:0]
		}
		msgs = append(msgs, msg)
	}
}

func CollectConsumer(consumer *KafkaConsumer, likeRankConsumer *interactive.KafkaConsumer) []events.Consumer {
	return []events.Consumer{consumer, likeRankConsumer}
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
		MaxBytes:       10e6,                   // 10MB 表示消费者可接受的最大批量大小, broker将截断消息以满足此最大值 比如MinBytes=10e3, MaxBytes=10e6, 则broker将返回10KB到10MB之间的消息
		MaxWait:        500 * time.Millisecond, // 500ms内有数据就返回, 即使没达到MinBytes, 与MinBytes取最小值, 有一个满足就返回
		CommitInterval: 0,                      // 多久自动commit一次offset 0表示同步提交
		StartOffset:    kafka.LastOffset,       // 从最新的offset开始读取
	})
	return r
}
