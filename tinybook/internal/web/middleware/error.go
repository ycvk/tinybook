package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorLogMiddlewareBuilder struct {
	log *zap.Logger
}

func NewErrorLogMiddleware(log *zap.Logger) *ErrorLogMiddlewareBuilder {
	return &ErrorLogMiddlewareBuilder{log: log}
}

func (builder *ErrorLogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			for _, e := range ctx.Errors {
				// 对每个错误进行日志记录
				builder.log.Error("请求出现错误", zap.Any("error", e.Err.Error()))
			}
		}
	}
}
