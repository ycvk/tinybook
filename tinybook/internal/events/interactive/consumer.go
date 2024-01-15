package interactive

import (
	"context"
	"github.com/Yiling-J/theine-go"
	"github.com/bytedance/sonic"
	"github.com/cockroachdb/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
	"tinybook/tinybook/internal/domain"
	"tinybook/tinybook/pkg/kafkax"
)

const (
	GroupLikeRankRead = "group-article-like-rank"
	LikeRankLocalFlag = "rank:like_rank_local_flag"
	RedisLikeRankKey  = "article:like_count"
	TopLikeRankNum    = 100
)

var (
	TimeToRefreshLocalCache = 1 * time.Minute  // 定时检查本地缓存是否需要更新
	TimeToCommitOffset      = 30 * time.Second // 多久自动commit一次offset
)

type KafkaConsumer struct {
	reader    *kafka.Reader
	log       *zap.Logger
	cli       *theine.Cache[string, any]
	redisCli  redis.Cmdable
	mu        sync.Mutex // 用于同步
	timer     *time.Timer
	isWaiting bool
}

func NewKafkaLikeRankConsumer(log *zap.Logger, cli *theine.Cache[string, any], redisCli redis.Cmdable) *KafkaConsumer {
	reader := InitReader(GroupLikeRankRead, TopicInteractiveLikeRank)
	collector := kafkax.NewReaderCollector(reader) // 用于收集 Kafka 读取器的统计信息
	prometheus.MustRegister(collector)             // 注册 Prometheus
	return &KafkaConsumer{
		log:       log,
		reader:    reader,
		cli:       cli,
		redisCli:  redisCli,
		timer:     nil,
		isWaiting: false,
	}
}

func (k *KafkaConsumer) Start() {
	ctx := context.Background()
	go func() {
		k.log.Info("like rank consumer start")
		k.Consume(ctx) // 消费kafka消息
	}()
	go func() {
		ctx := context.Background()
		k.log.Info("like rank consumer ticker start")
		k.Ticker(ctx, TimeToRefreshLocalCache) // 定时检查本地缓存是否需要更新
	}()
}

func (k *KafkaConsumer) Consume(ctx context.Context) {
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		k.log.Info("close kafka like rank consumer")
		if err != nil {
			k.log.Error("close kafka like rank consumer failed", zap.Error(err))
		}
	}(k.reader)
	for {
		// 读取消息
		message, err := k.reader.FetchMessage(ctx)
		if err != nil {
			k.log.Error("read message failed", zap.Error(err))
			continue
		}
		// 解析消息
		var event LikeRankEvent
		err = sonic.Unmarshal(message.Value, &event)
		if err != nil {
			k.log.Error("like rank consumer unmarshal message failed", zap.Error(err))
			continue
		}
		// 设置redis缓存更新标志位
		if event.Change {
			k.Call(func() { k.redisCli.Set(ctx, LikeRankLocalFlag, 1, 0) }, TimeToCommitOffset)
		}
		// 提交offset
		err = k.reader.CommitMessages(ctx, message)
	}
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
		CommitInterval: TimeToCommitOffset,     // 多久自动commit一次offset 0表示同步提交
		StartOffset:    kafka.LastOffset,       // 从最新的offset开始读取
	})
	return r
}

func (k *KafkaConsumer) Ticker(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			goto tickerEnd
		case <-ticker.C:
			// 每固定时间检查一次本地缓存是否需要更新
			change, err := k.redisCli.Get(ctx, LikeRankLocalFlag).Bool()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					// redis缓存标志位不存在, 说明没有新的点赞数更新
				} else {
					k.log.Error("ticker get like rank local flag failed", zap.Error(err))
				}
				continue
			}
			if change {
				// redis缓存标志位存在, 说明有新的点赞数更新
				// 从redis中获取 topN 文章的点赞数与id
				topNLike, err := k.redisCli.ZRevRangeWithScores(ctx, RedisLikeRankKey, 0, TopLikeRankNum-1).Result()
				if err != nil {
					k.log.Error("ticker get topN like rank from redis failed", zap.Error(err))
					continue
				}
				// 获取到的 redis.Z 转换为 domain.Interactive
				interactivesMap := lo.Map(topNLike, func(item redis.Z, index int) domain.Interactive {
					s := item.Member.(string)
					id, _ := strconv.ParseInt(s, 10, 64)
					return domain.Interactive{
						BizId:     id,
						LikeCount: int64(item.Score),
					}
				})
				// 将topNLike写入本地缓存
				k.log.Info("ticker set topN like rank to local cache")
				k.cli.Set(RedisLikeRankKey, interactivesMap, 1)
				// 删除redis缓存标志位
				k.redisCli.Del(ctx, LikeRankLocalFlag)
			}
		}
	}
tickerEnd:
	k.log.Info("consumer ticker end")
	return
}

// Call 执行函数，但保证在给定时间内最多执行一次
func (k *KafkaConsumer) Call(f func(), duration time.Duration) {
	k.mu.Lock()

	if k.isWaiting {
		k.mu.Unlock()
		return
	}
	k.isWaiting = true
	k.mu.Unlock()

	go func() {
		f()
		time.Sleep(duration)
		k.mu.Lock()
		k.isWaiting = false
		k.mu.Unlock()
	}()
}
