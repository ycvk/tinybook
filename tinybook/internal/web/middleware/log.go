package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type LogMiddlewareBuilder struct {
	logFn              func(ctx *gin.Context, accessLog *AccessLog)
	allowPrintReqBody  bool
	allowPrintRespBody bool
}

type AccessLog struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Ip       string `json:"ip"`
	ReqBody  string `json:"req_body"`
	RespBody string `json:"resp_body"`
	Duration string `json:"duration"`
	Status   int    `json:"status"`
}

func NewLogMiddleware(logFn func(ctx *gin.Context, accessLog *AccessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFn:              logFn,
		allowPrintReqBody:  false,
		allowPrintRespBody: false,
	}
}

func (builder *LogMiddlewareBuilder) AllowPrintReqBody() *LogMiddlewareBuilder {
	builder.allowPrintReqBody = true
	return builder
}

func (builder *LogMiddlewareBuilder) AllowPrintRespBody() *LogMiddlewareBuilder {
	builder.allowPrintRespBody = true
	return builder
}

func (builder *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Request.URL.Path
		if len(url) > 1024 {
			url = url[:1024]
		}
		method := ctx.Request.Method
		ip := ctx.ClientIP()
		accessLog := &AccessLog{
			Path:   url,
			Method: method,
			Ip:     ip,
		}
		if builder.allowPrintReqBody {
			data, _ := ctx.GetRawData()
			if len(data) > 2048 {
				accessLog.ReqBody = string(data[:2048])
			} else {
				accessLog.ReqBody = string(data)
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		}
		now := time.Now()
		if builder.allowPrintRespBody {
			rw := &responseWriter{ResponseWriter: ctx.Writer, al: accessLog}
			ctx.Writer = rw
		}
		defer func() {
			accessLog.Duration = time.Since(now).String()
			if builder.logFn != nil {
				builder.logFn(ctx, accessLog)
			}
		}()
		// 执行下一个中间件
		ctx.Next()
	}
}

// responseWriter 重写gin.ResponseWriter 用于获取响应body
type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.al.RespBody = string(data)
	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) WriteString(s string) (int, error) {
	return rw.ResponseWriter.WriteString(s)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.al.Status = code
	rw.ResponseWriter.WriteHeader(code)
}
