package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"strings"
	"time"
	web2 "tinybook/tinybook/article/web"
	"tinybook/tinybook/internal/web"
	"tinybook/tinybook/internal/web/jwt"
	"tinybook/tinybook/internal/web/middleware"
	"tinybook/tinybook/pkg/ginx/middleware/prometheus"
	"tinybook/tinybook/pkg/ginx/middleware/ratelimit"
	"tinybook/tinybook/pkg/limiter"
)

func InitWebServer(handlerFunc []gin.HandlerFunc, userHandler *web.UserHandler,
	wechatHandler *web.OAuth2WechatHandler, articleHandler *web2.ArticleHandler) *gin.Engine {
	engine := gin.Default()
	// 注册中间件
	engine.Use(handlerFunc...)
	// 接入opentelemetry
	g := otelgin.Middleware("tinybook")
	engine.Use(g)
	// 注册用户路由
	userHandler.RegisterRoutes(engine)
	// 注册文章路由
	articleHandler.RegisterRoutes(engine)
	// 注册wechat oauth2路由
	wechatHandler.RegisterRoutes(engine)
	return engine
}

// InitHandlerFunc 初始化中间件
func InitHandlerFunc(redisClient redis.Cmdable, handler jwt.Handler, logger *zap.Logger) []gin.HandlerFunc {
	corsConfig := initCorsConfig() // 跨域配置
	//rateLimit := initRateLimit(redisClient)        // 限流器
	log := initLogger(logger)                       // 日志请求记录器
	errorLog := initErrorLog(logger)                // 错误日志
	loginJWT := initLoginJWT(handler)               // 登录jwt
	prometheusRespTime := initPrometheusRespTime()  // prometheus响应时间
	activeReq := initPrometheusActiveReq()          // prometheus活跃链接数
	middleware.InitCounter(prometheus2.CounterOpts{ // prometheus请求code
		Namespace: "tinybook",
		Subsystem: "gin",
		Name:      "req_code",
		Help:      "统计gin的http接口请求code",
	})

	return []gin.HandlerFunc{
		corsConfig,
		//rateLimit,
		log,
		errorLog,
		loginJWT,
		prometheusRespTime,
		activeReq,
	}
}

// initCorsConfig 跨域配置
func initCorsConfig() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},         //允许跨域的方法
		AllowHeaders:     []string{"Content-Type", "Authorization"},  // 允许跨域的Header
		ExposeHeaders:    []string{"X-Jwt-Token", "X-Refresh-Token"}, // 允许访问的响应头
		AllowCredentials: true,                                       //  允许携带cookie
		AllowOriginFunc: func(origin string) bool { //允许跨域的域名
			return strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour, //缓存时间
	})
}

func initErrorLog(logger *zap.Logger) gin.HandlerFunc {
	return middleware.NewErrorLogMiddleware(logger).Build()
}

func initPrometheusRespTime() gin.HandlerFunc {
	return prometheus.NewBuilder("tinybook", "gin", "http_request", "统计gin的http接口请求数据", "1").BuildResponseTime()
}

func initPrometheusActiveReq() gin.HandlerFunc {
	return prometheus.NewBuilder("tinybook", "gin", "http_request", "统计gin的活跃链接数", "1").BuildActiveRequest()
}

// initLogger 初始化日志请求记录器
func initLogger(logger *zap.Logger) gin.HandlerFunc {
	return middleware.NewLogMiddleware(func(ctx *gin.Context, accessLog *middleware.AccessLog) {
		logger.Info("HTTP Log",
			zap.String("Path", accessLog.Path),
			zap.String("Method", accessLog.Method),
			zap.String("IP", accessLog.Ip),
			zap.String("RequestBody", accessLog.ReqBody),
			zap.String("ResponseBody", accessLog.RespBody),
			zap.String("Duration", accessLog.Duration),
			zap.Int("Status", accessLog.Status),
		)
	}).AllowPrintReqBody().AllowPrintRespBody().Build() // 允许打印请求和响应的body
}

// initLoginJWT 初始化登录jwt
func initLoginJWT(handler jwt.Handler) gin.HandlerFunc {
	middlewareBuilder := middleware.NewLoginJWTMiddlewareBuilder(handler)
	return middlewareBuilder.Build()
}

// initRateLimit 初始化限流器
func initRateLimit(redisClient redis.Cmdable) gin.HandlerFunc {
	return ratelimit.NewBuilder(limiter.NewRedisSlideWindowLimiter(redisClient, time.Second, 5)).Build() // 一秒钟限制5次
}

// initLoginSession 初始化登录session
func initLoginSession(engine *gin.Engine) {
	type Config struct {
		Host string `yaml:"addr"`
	}
	var cfg Config
	err2 := viper.UnmarshalKey("redis", &cfg)
	if err2 != nil {
		panic(err2)
	}
	loginMiddleware := middleware.LoginMiddlewareBuilder{}
	store, err := redisSession.NewStore(16,
		"tcp",
		cfg.Host,
		"",
		[]byte("zcPbUOs7zYO1ky2WgE14chotKwcp95Hp"), //authentication key 身份验证密钥
		[]byte("GdGvU8pRs439iNREpNtl1gZhY7jU8zRt"), //encryption key 加密密钥
	)
	if err != nil {
		panic(err)
	}
	engine.Use(
		loginMiddleware.Build(),          // 初始化登录中间件
		sessions.Sessions("ssid", store), // 初始化session
	)
}
