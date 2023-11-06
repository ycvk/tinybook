package ioc

import (
	"geek_homework/tinybook/config"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/middleware"
	"geek_homework/tinybook/pkg/ginx/middleware/ratelimit"
	"geek_homework/tinybook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(handlerFunc []gin.HandlerFunc, userHandler *web.UserHandler, wechatHandler *web.OAuth2WechatHandler) *gin.Engine {
	engine := gin.Default()
	// 注册中间件
	engine.Use(handlerFunc...)
	// 注册用户路由
	userHandler.RegisterRoutes(engine)
	// 注册oauth2路由
	wechatHandler.RegisterRoutes(engine)
	return engine
}

// InitHandlerFunc 初始化中间件
func InitHandlerFunc(redisClient redis.Cmdable) []gin.HandlerFunc {
	corsConfig := initCorsConfig()          // 跨域配置
	rateLimit := initRateLimit(redisClient) // 限流器
	loginJWT := initLoginJWT()              // 登录jwt
	return []gin.HandlerFunc{corsConfig, rateLimit, loginJWT}
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

// initLoginJWT 初始化登录jwt
func initLoginJWT() gin.HandlerFunc {
	middlewareBuilder := middleware.LoginJWTMiddlewareBuilder{}
	return middlewareBuilder.Build()
}

// initRateLimit 初始化限流器
func initRateLimit(redisClient redis.Cmdable) gin.HandlerFunc {
	return ratelimit.NewBuilder(limiter.NewRedisSlideWindowLimiter(redisClient, time.Second, 5)).Build() // 一秒钟限制5次
}

// initLoginSession 初始化登录session
func initLoginSession(engine *gin.Engine) {
	loginMiddleware := middleware.LoginMiddlewareBuilder{}
	store, err := redisSession.NewStore(16,
		"tcp",
		config.Config.Redis.Host,
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
