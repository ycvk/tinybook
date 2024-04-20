package cache

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"time"
	"tinybook/tinybook/internal/domain"
)

var ErrUserNotFound = redis.Nil

type UserCache interface {
	GetById(ctx context.Context, id int64) (domain.User, error)
	SetById(ctx context.Context, user domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c RedisUserCache) GetById(ctx context.Context, id int64) (domain.User, error) {
	key := c.key(id)
	result, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	// 反序列化
	err = sonic.UnmarshalString(result, &user)
	return user, err
}

// key 生成缓存key
func (c RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (c RedisUserCache) SetById(ctx context.Context, user domain.User) error {
	key := c.key(user.Id)
	// 序列化
	value, err := sonic.MarshalString(user)
	if err != nil {
		return err
	}
	// 设置缓存
	return c.cmd.Set(ctx, key, value, c.expiration).Err()
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
