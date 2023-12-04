//go:build wireinject

package main

import (
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/jwt"
	"geek_homework/tinybook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {

	wire.Build(
		// 初始化redis, db, localCache, mongoDB
		ioc.InitRedis, ioc.InitDB, ioc.InitLocalCache, ioc.InitMongoDB, ioc.InitMongoDBV2,
		// 初始化user模块
		cache.NewRedisUserCache, dao.NewGormUserDAO, repository.NewCachedUserRepository, service.NewUserService,
		// 初始化code模块
		cache.NewLocalCodeCache, repository.NewCachedCodeRepository, service.NewCodeService,
		// 初始化sms模块
		ioc.InitSMSService, repository.NewGormSMSRepository, dao.NewGormSMSDAO,
		// 初始化article模块
		repository.NewCachedArticleRepository, dao.NewMongoDBArticleDAO, service.NewArticleService, cache.NewRedisArticleCache,
		// 初始化interactive模块
		cache.NewRedisInteractiveCache, repository.NewCachedInteractiveRepository, dao.NewGormInteractiveDAO, service.NewInteractiveService,
		// 初始化oauth2模块
		ioc.InitWechatService,
		// 初始化handler
		web.NewUserHandler, web.NewOAuth2WechatHandler, jwt.NewRedisJWTHandler,
		web.NewArticleHandler,
		// 初始化web 和 中间件
		ioc.InitWebServer, ioc.InitHandlerFunc, ioc.InitLogger,
	)

	return gin.Default()
}
