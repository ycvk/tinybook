package main

import (
	"geek_homework/tinybook/config"
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/middleware"
	"geek_homework/tinybook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	engine := gin.Default()
	// 跨域配置
	initCorsConfig(engine)
	// 初始化限流
	redisClient := initRedis()
	build := ratelimit.NewBuilder(redisClient, time.Second, 5).Build() // 一秒钟限制5次
	engine.Use(build)
	// 初始化登录session
	//initLoginSession(engine)
	// 初始化登录jwt
	initLoginJWT(engine)
	// 初始化数据库
	db := initDB()
	// 初始化用户模块
	initUser(db, engine)

	//engine.GET("/ping", func(ctx *gin.Context) {
	//	ctx.String(200, "hello world!")
	//})
	engine.Run(":8081")
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Host,
	})
}

// initLoginJWT 初始化登录jwt
func initLoginJWT(engine *gin.Engine) {
	middlewareBuilder := middleware.LoginJWTMiddlewareBuilder{}
	engine.Use(middlewareBuilder.Build())
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

func initUser(db *gorm.DB, engine *gin.Engine) {
	userDAO := dao.NewUserDAO(db)
	userRepository := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)

	userHandler.RegisterRoutes(engine)
}

// initDB 初始化数据库
func initDB() *gorm.DB {
	dsn := config.Config.DB.Host
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()
	if err != nil {
		panic(err)
	}
	return db
}

// initCorsConfig 跨域配置
func initCorsConfig(engine *gin.Engine) {
	corsConfig := cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},        //允许跨域的方法
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // 允许跨域的Header
		ExposeHeaders:    []string{"X-Jwt-Token"},                   // 允许访问的响应头
		AllowCredentials: true,                                      //  允许携带cookie
		AllowOriginFunc: func(origin string) bool { //允许跨域的域名
			return strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour, //缓存时间
	})
	engine.Use(corsConfig)
}
