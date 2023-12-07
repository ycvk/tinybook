package ioc

import (
	"geek_homework/tinybook/internal/events"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/jwt"
	"geek_homework/tinybook/internal/web/middleware"
	"geek_homework/tinybook/pkg/ginx/middleware/ratelimit"
	"geek_homework/tinybook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"time"
)

func InitWebServer(handlerFunc []gin.HandlerFunc, userHandler *web.UserHandler,
	wechatHandler *web.OAuth2WechatHandler, articleHandler *web.ArticleHandler, consumers []events.Consumer) *gin.Engine {
	engine := gin.Default()
	// 注册中间件
	engine.Use(handlerFunc...)
	// 注册用户路由
	userHandler.RegisterRoutes(engine)
	// 注册文章路由
	articleHandler.RegisterRoutes(engine)
	// 注册wechat oauth2路由
	wechatHandler.RegisterRoutes(engine)
	// 注册kafka消费者
	for _, consumer := range consumers {
		consumer.Start()
	}
	return engine
}

// InitHandlerFunc 初始化中间件
func InitHandlerFunc(redisClient redis.Cmdable, handler jwt.Handler, logger *zap.Logger) []gin.HandlerFunc {
	corsConfig := initCorsConfig()          // 跨域配置
	rateLimit := initRateLimit(redisClient) // 限流器
	log := initLogger(logger)               // 日志 todo 本地开发时，可以注释掉
	errorLog := initErrorLog(logger)        // 错误日志
	loginJWT := initLoginJWT(handler)       // 登录jwt
	return []gin.HandlerFunc{corsConfig, rateLimit, log, errorLog, loginJWT}
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

// initLogger 初始化日志
func initLogger(logger *zap.Logger) gin.HandlerFunc {
	return middleware.NewLogMiddleware(func(ctx *gin.Context, accessLog *middleware.AccessLog) {
		logger.Info("", zap.Any("accessLog", accessLog))
	}).AllowPrintReqBody().AllowPrintRespBody().Build()
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
