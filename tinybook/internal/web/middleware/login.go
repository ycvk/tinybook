package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

// Build 构建登录中间件
func (builder *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 注册time.Time类型, 否则session无法保存time.Time类型
	gob.Register(time.Time{})

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// 登录和注册不需要经过登录中间件, 可以在未经认证的情况下访问
		if path == "/users/login" || path == "/users/signup" {
			return
		}
		session := sessions.Default(ctx)
		userId := session.Get("userId")
		if userId == nil {
			// 401未授权的状态码中止当前请求
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// session刷新登录态
		const updateTimeKey = "update_time"
		now := time.Now()
		updateTime := session.Get(updateTimeKey)
		lastUpdateTime, ok := updateTime.(time.Time)      // 类型断言
		if !ok || now.Sub(lastUpdateTime) > time.Minute { // 一分钟刷新一次
			session.Set(updateTimeKey, now)
			session.Set("userId", userId)
			err := session.Save()
			if err != nil {
				slog.Info("session刷新失败", "err", err)
			}
		}
	}
}
