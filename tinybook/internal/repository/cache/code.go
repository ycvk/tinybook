package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

var ErrCodeVerifyTooMany = errors.New("验证码错误次数过多, 请稍后重试")

type CodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) *CodeCache {
	return &CodeCache{
		cmd: cmd,
	}
}

// SetCode 设置验证码 timeInterval: 有效时间, 比如600 表示10分钟内有效
func (c *CodeCache) SetCode(ctx context.Context, phone, biz, code, timeInterval string) error {
	i, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code, timeInterval).Int()
	if err != nil {
		return err
	}
	switch i {
	case -1:
		slog.Error("key没有设置过期时间 ", "key", c.key(biz, phone))
		return errors.New("验证码没有设置过期时间")
	case 1:
		slog.Error("发送验证码太频繁 ", "key", c.key(biz, phone))
		return errors.New("发送验证码太频繁")
	default:
		return nil
	}
}

func (c *CodeCache) VerifyCode(ctx context.Context, phone, biz, code string) (bool, error) {
	i, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
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

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("code:%s:%s", biz, phone)
}
