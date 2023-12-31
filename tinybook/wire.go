//go:build wireinject

package main

import (
	"geek_homework/tinybook/internal/events/article"
	"geek_homework/tinybook/internal/events/interactive"
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/jwt"
	"geek_homework/tinybook/ioc"
	"github.com/google/wire"
)

// 热榜服务
var rankingServiceProvider = wire.NewSet(
	cache.NewRedisRankingCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

// interactive 互动服务
var interactiveServiceProvider = wire.NewSet(
	cache.NewRedisInteractiveCache,
	dao.NewGormInteractiveDAO,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

func InitWebServer() *App {

	wire.Build(
		// 初始化redis, db, localCache, mongoDB
		ioc.InitRedis, ioc.InitDB, ioc.InitLocalCache, ioc.InitMongoDB, ioc.InitMongoDBV2,
		// 初始化redisLock
		ioc.InitRedisLock,
		// 初始化user模块
		cache.NewRedisUserCache, dao.NewGormUserDAO, repository.NewCachedUserRepository, service.NewUserService,
		// 初始化code模块
		cache.NewLocalCodeCache, repository.NewCachedCodeRepository, service.NewCodeService,
		// 初始化sms模块
		ioc.InitSMSService, repository.NewGormSMSRepository, dao.NewGormSMSDAO,
		// 初始化article模块
		repository.NewCachedArticleRepository, dao.NewMongoDBArticleDAO, service.NewArticleService, cache.NewRedisArticleCache,
		// 初始化interactive模块
		interactiveServiceProvider,
		// 初始化oauth2模块
		ioc.InitWechatService,
		// 初始化ranking模块
		rankingServiceProvider, ioc.InitJobs, ioc.InitRankingJob,
		// 初始化handler
		web.NewUserHandler, web.NewOAuth2WechatHandler, jwt.NewRedisJWTHandler,
		web.NewArticleHandler,
		// 初始化web 和 中间件
		ioc.InitWebServer, ioc.InitHandlerFunc, ioc.InitLogger,
		// 初始化阅读数 read num kafka
		ioc.InitWriter, article.NewKafkaArticleProducer,
		article.NewKafkaConsumer, article.CollectConsumer,
		// 初始化点赞榜 like rank kafka
		interactive.NewKafkaLikeRankProducer, interactive.NewKafkaLikeRankConsumer,

		wire.Struct(new(App), "*"),
	)

	return &App{}
}
