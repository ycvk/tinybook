package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
}

// Build 构建登录中间件
func (builder *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// 登录和注册不需要经过登录中间件, 可以在未经认证的情况下访问
		if path == "/users/login" || path == "/users/signup" {
			return
		}
		session := sessions.Default(ctx)
		if session.Get("userId") == nil {
			// 401未授权的状态码中止当前请求
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
