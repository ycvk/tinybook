package article

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
)

const TopicArticleRead = "topic-article-read"

type ReadEvent struct {
	ArticleID int64 `json:"article_id"`
	UserID    int64 `json:"user_id"`
}

type ReadEventProducer interface {
	ProduceReadEvent(event ReadEvent) error
}

type KafkaAsyncProducer struct {
	writer *kafka.Writer
}

func NewKafkaAsyncProducer(writer *kafka.Writer) ReadEventProducer {
	return &KafkaAsyncProducer{writer: writer}
}

func (k *KafkaAsyncProducer) ProduceReadEvent(event ReadEvent) error {
	bytes, err := sonic.Marshal(event)
	if err != nil {
		return err
	}
	err = k.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: TopicArticleRead,
		Value: bytes,
	})
	return err
}
