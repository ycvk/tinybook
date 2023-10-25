package ioc

import (
	"geek_homework/tinybook/config"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/middleware"
	"geek_homework/tinybook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(handlerFunc []gin.HandlerFunc, handler *web.UserHandler) *gin.Engine {
	engine := gin.Default()
	// 注册中间件
	engine.Use(handlerFunc...)
	// 注册路由
	handler.RegisterRoutes(engine)
	return engine
}

func InitHandlerFunc(redisClient redis.Cmdable) []gin.HandlerFunc {
	corsConfig := initCorsConfig()
	rateLimit := initRateLimit(redisClient)
	loginJWT := initLoginJWT()
	return []gin.HandlerFunc{corsConfig, rateLimit, loginJWT}
}

// initCorsConfig 跨域配置
func initCorsConfig() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},        //允许跨域的方法
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // 允许跨域的Header
		ExposeHeaders:    []string{"X-Jwt-Token"},                   // 允许访问的响应头
		AllowCredentials: true,                                      //  允许携带cookie
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

func initRateLimit(redisClient redis.Cmdable) gin.HandlerFunc {
	return ratelimit.NewBuilder(redisClient, time.Second, 5).Build() // 一秒钟限制5次

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
