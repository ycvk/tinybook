package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/url"
	"strconv"
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

var vector *prometheus.CounterVec

func InitCounter(opts prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opts, []string{"status"})
	prometheus.MustRegister(vector)
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
		urlPath := ctx.Request.URL.Path
		if len(urlPath) > 1024 {
			urlPath = urlPath[:1024]
		}
		method := ctx.Request.Method
		ip := ctx.ClientIP()
		accessLog := &AccessLog{
			Path:   urlPath,
			Method: method,
			Ip:     ip,
		}
		if builder.allowPrintReqBody {
			data, _ := ctx.GetRawData()
			var unquote []byte
			if len(data) > 2048 {
				unquote, _ = unicodeToZH(data[:2048])
			} else {
				unquote, _ = unicodeToZH(data)
			}
			accessLog.ReqBody = string(unquote)
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
			vector.WithLabelValues(strconv.Itoa(accessLog.Status)).Inc() //prometheus 记录状态码
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
	zh, _ := unicodeToZH(data)
	rw.al.RespBody = string(zh)
	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) WriteString(s string) (int, error) {
	return rw.ResponseWriter.WriteString(s)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.al.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

// unicodeToZH unicode转中文
func unicodeToZH(raw []byte) ([]byte, error) {
	//str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	//if err != nil {
	//	return nil, err
	//}
	//return []byte(str), nil
	// 用url.QueryUnescape替换strconv.Unquote 更高效
	decodedStr, err := url.QueryUnescape(string(raw))
	if err != nil {
		return nil, err
	}
	return []byte(decodedStr), nil
}
