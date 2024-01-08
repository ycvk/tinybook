package ioc

import (
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
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: cfg.Addr,
		})
	})
	return redisClient
}
