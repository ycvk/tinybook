package ioc

import (
	"geek_homework/tinybook/config"
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	Once        sync.Once
	redisClient redis.Cmdable
)

func InitRedis() redis.Cmdable {
	Once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: config.Config.Redis.Host,
		})
	})
	return redisClient
}
