package interactive

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
)

const TopicInteractiveLikeRank = "topic-article-like-rank"

type LikeRankEvent struct {
	ArticleID int64 `json:"article_id"`
	LikeCount int64 `json:"like_count"`
	Change    bool  `json:"change"`
}

type LikeRankEventProducer interface {
	ProduceLikeRankEvent(event LikeRankEvent) error
}

type KafkaLikeRankProducer struct {
	writer *kafka.Writer
}

func NewKafkaLikeRankProducer(writer *kafka.Writer) LikeRankEventProducer {
	return &KafkaLikeRankProducer{writer: writer}
}

func (k KafkaLikeRankProducer) ProduceLikeRankEvent(event LikeRankEvent) error {
	bytes, err := sonic.Marshal(event)
	if err != nil {
		return err
	}
	err = k.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: TopicInteractiveLikeRank,
		Value: bytes,
	})
	return err
}
