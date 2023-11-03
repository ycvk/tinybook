//go:build wireinject

package main

import (
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {

	wire.Build(
		// 初始化redis 和 db 和 localCache
		ioc.InitRedis, ioc.InitDB, ioc.InitLocalCache,
		// 初始化user模块
		cache.NewRedisUserCache, dao.NewGormUserDAO, repository.NewCachedUserRepository, service.NewUserService,
		// 初始化code模块
		cache.NewLocalCodeCache, repository.NewCachedCodeRepository, service.NewCodeService, dao.NewGormCodeDAO,
		// 初始化sms模块
		ioc.InitSMSService,
		// 初始化handler
		web.NewUserHandler,
		// 初始化web
		ioc.InitWebServer, ioc.InitHandlerFunc,
	)

	return gin.Default()
}
