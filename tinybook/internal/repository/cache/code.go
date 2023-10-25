package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/Yiling-J/theine-go"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
	"time"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

var ErrCodeVerifyTooMany = errors.New("验证码错误次数过多, 请稍后重试")

type CodeCache interface {
	SetCode(ctx context.Context, phone, biz, code, timeInterval string) error
	VerifyCode(ctx context.Context, phone, biz, code string) (bool, error)
}

type LocalCodeCache struct {
	client *theine.Cache[string, any]
}

func (l *LocalCodeCache) SetCode(ctx context.Context, phone, biz, code, timeInterval string) error {
	atoi, err := strconv.Atoi(timeInterval)
	if err != nil {
		return err
	}
	s := key(biz, phone)
	// 检查是否可以发送下次验证码
	flag, b := l.client.Get(s + "-ttl")
	if b && flag.(bool) {
		//无论key存在还是为true, 都不允许发送验证码
		return errors.New("验证码发送太频繁, 请60秒后重试")
	}

	// 设置最大验证次数
	l.client.SetWithTTL(s+"-limit", 3, 1, time.Second*time.Duration(atoi))
	// 设置验证码
	ttl := l.client.SetWithTTL(s, code, 1, time.Second*time.Duration(atoi))
	// 设置是否可以发送下次验证码 60秒内不允许发送
	l.client.SetWithTTL(s+"-ttl", true, 1, time.Second*60)
	if !ttl {
		slog.Error("设置本地缓存失败 ", "key", s)
		return errors.New("设置本地缓存失败")
	}
	return nil
}

func (l *LocalCodeCache) VerifyCode(ctx context.Context, phone, biz, code string) (bool, error) {
	key := key(biz, phone)
	// 检查是否超过次数
	cnt, ok := l.client.Get(key + "-limit")
	if !ok {
		return false, errors.New("验证码已过期或不存在")
	}
	if cnt.(int) <= 0 {
		return false, ErrCodeVerifyTooMany
	}
	l.client.Set(key+"-limit", cnt.(int)-1, 1)

	get, b := l.client.Get(key)
	if !b {
		return false, errors.New("验证码已过期或不存在")
	}
	if get.(string) == code {
		l.client.Delete(key)
		return true, nil
	}
	return false, errors.New("验证码错误, 请重试")
}

func NewLocalCodeCache(cache *theine.Cache[string, any]) CodeCache {
	return &LocalCodeCache{
		client: cache,
	}
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

// SetCode 设置验证码 timeInterval: 有效时间, 比如600 表示10分钟内有效
func (c *RedisCodeCache) SetCode(ctx context.Context, phone, biz, code, timeInterval string) error {
	i, err := c.cmd.Eval(ctx, luaSetCode, []string{key(biz, phone)}, code, timeInterval).Int()
	if err != nil {
		return err
	}
	switch i {
	case -1:
		slog.Error("key没有设置过期时间 ", "key", key(biz, phone))
		return errors.New("验证码没有设置过期时间")
	case 1:
		slog.Error("发送验证码太频繁 ", "key", key(biz, phone))
		return errors.New("发送验证码太频繁")
	default:
		return nil
	}
}

func (c *RedisCodeCache) VerifyCode(ctx context.Context, phone, biz, code string) (bool, error) {
	i, err := c.cmd.Eval(ctx, luaVerifyCode, []string{key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch i {
	case 0:
		return true, nil
	case 1:
		return false, errors.New("验证码错误, 请重试")
	case -2:
		return false, ErrCodeVerifyTooMany
	default:
		return false, errors.New("验证码已过期或不存在")
	}
}

func key(biz, phone string) string {
	return fmt.Sprintf("code:%s:%s", biz, phone)
}
