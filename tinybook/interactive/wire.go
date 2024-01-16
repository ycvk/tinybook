//go:build wireinject

package main

import (
	"github.com/google/wire"
	"tinybook/tinybook/interactive/events/rank"
	"tinybook/tinybook/interactive/events/readcount"
	"tinybook/tinybook/interactive/grpc"
	"tinybook/tinybook/interactive/ioc"
	"tinybook/tinybook/interactive/repository"
	"tinybook/tinybook/interactive/repository/cache"
	"tinybook/tinybook/interactive/repository/dao"
	"tinybook/tinybook/interactive/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitRedis, ioc.InitLogger, ioc.InitLocalCache,
)

var interactiveServiceSet = wire.NewSet(
	dao.NewGormInteractiveDAO, cache.NewRedisInteractiveCache,
	repository.NewCachedInteractiveRepository, service.NewInteractiveService,
)

func InitInteractiveApp() *App {
	wire.Build(
		thirdPartySet,
		interactiveServiceSet,
		// 初始化 kafka writer
		ioc.InitWriter,
		// 初始化阅读数消费者 read num kafka
		readcount.NewKafkaReadCountConsumer,
		// 初始化点赞榜 like rank kafka
		rank.NewKafkaLikeRankProducer, rank.NewKafkaLikeRankConsumer,
		// 收集所有的consumer
		readcount.CollectConsumer,
		// 初始化 grpc server
		grpc.NewInteractiveServiceServer, ioc.InitGrpcServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
