package main

import (
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	engine := gin.Default()
	// 初始化登录session
	loginMiddleware := middleware.LoginMiddlewareBuilder{}
	store := cookie.NewStore([]byte("secret"))
	// 跨域配置
	engine.Use(initCorsConfig(),
		sessions.Sessions("ssid", store),
		loginMiddleware.Build())
	// 初始化数据库
	db := initDB()
	// 初始化用户模块
	initUser(db, engine)

	engine.Run(":8080")
}

func initUser(db *gorm.DB, engine *gin.Engine) {
	userDAO := dao.NewUserDAO(db)
	userRepository := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)

	userHandler.RegisterRoutes(engine)
}

func initDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/ycvk?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()
	if err != nil {
		panic(err)
	}
	return db
}

func initCorsConfig() gin.HandlerFunc {
	corsConfig := cors.New(cors.Config{
		//AllowMethods:     []string{"POST", "GET", "OPTIONS"},                   //允许跨域的方法
		AllowHeaders: []string{"Content-Type", "Authorization"}, // 允许跨域的Header
		//ExposeHeaders:    []string{"Content-Length"},                           // 允许访问的Header
		AllowCredentials: true, //  允许携带cookie
		AllowOriginFunc: func(origin string) bool { //允许跨域的域名
			return strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour, //缓存时间
	})
	return corsConfig
}
