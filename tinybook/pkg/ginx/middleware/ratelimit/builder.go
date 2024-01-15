package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tinybook/tinybook/pkg/limiter"
)

type Builder struct {
	prefix  string
	limiter limiter.Limiter
}

func NewBuilder(limiter2 limiter.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: limiter2,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
		limited, err := b.limiter.Limit(ctx, key)
		if err != nil {
			slog.Error("Rate limit failed", "err", err)
			// 如果这边出错了
			// 保守做法：因为借助于 Redis 来做限流，那么 Redis 崩溃了，为了防止系统崩溃，直接限流
			ctx.AbortWithStatus(http.StatusInternalServerError)
			// 激进做法：虽然 Redis 崩溃了，但是这个时候还是要尽量服务正常的用户，所以不限流
			// ctx.Next()
			return
		}
		if limited {
			slog.Error("Rate limit reached, too many requests", "ip", ctx.ClientIP())
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
