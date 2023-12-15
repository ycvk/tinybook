package ioc

import (
	"geek_homework/tinybook/pkg/redisx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"sync"
)

var (
	redisOnce   sync.Once
	redisClient redis.Cmdable
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}
	hook := redisx.NewPrometheusHook(prometheus.SummaryOpts{
		Namespace: "tinybook",
		Subsystem: "redis",
		Name:      "redis",
		Help:      "统计redis操作耗时",
	})
	redisOnce.Do(func() {
		//redisClient = redis.NewClient(&redis.Options{
		//	Addr: cfg.Addr,
		//})
		newClient := redis.NewClient(&redis.Options{
			Addr: cfg.Addr,
		})
		newClient.AddHook(hook)
		redisClient = newClient
	})

	return redisClient
}
