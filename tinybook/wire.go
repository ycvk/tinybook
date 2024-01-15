//go:build wireinject

package main

import (
	"github.com/google/wire"
	"tinybook/tinybook/interactive/events/rank"
	"tinybook/tinybook/interactive/events/readcount"
	repository2 "tinybook/tinybook/interactive/repository"
	cache2 "tinybook/tinybook/interactive/repository/cache"
	dao2 "tinybook/tinybook/interactive/repository/dao"
	service2 "tinybook/tinybook/interactive/service"
	"tinybook/tinybook/internal/job"
	"tinybook/tinybook/internal/repository"
	"tinybook/tinybook/internal/repository/cache"
	"tinybook/tinybook/internal/repository/dao"
	"tinybook/tinybook/internal/service"
	"tinybook/tinybook/internal/web"
	"tinybook/tinybook/internal/web/jwt"
	"tinybook/tinybook/ioc"
)

// 热榜服务
var rankingServiceProvider = wire.NewSet(
	cache.NewRedisRankingCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

// interactive 互动服务
var interactiveServiceProvider = wire.NewSet(
	cache2.NewRedisInteractiveCache,
	dao2.NewGormInteractiveDAO,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

// job 服务
var jobServiceProvider = wire.NewSet(
	service.NewCronJobService,
	repository.NewCronJobRepository,
	dao.NewGormCronJobDao,

	job.NewScheduler,
	job.NewLocalFuncExecutor,
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
		// 初始化kafka writer
		ioc.InitWriter,
		// 初始化阅读数 read num kafka
		readcount.NewKafkaReadCountProducer, readcount.NewKafkaReadCountConsumer,
		// 初始化点赞榜 like rank kafka
		rank.NewKafkaLikeRankProducer, rank.NewKafkaLikeRankConsumer,
		// 收集所有的consumer
		readcount.CollectConsumer,

		// 初始化job
		jobServiceProvider,
		wire.Struct(new(App), "*"),
	)

	return &App{}
}
