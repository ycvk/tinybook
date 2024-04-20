package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpc2 "tinybook/tinybook/interactive/grpc"
	"tinybook/tinybook/pkg/grpcx"
)

func InitGrpcServer(interactiveServer *grpc2.InteractiveServiceServer, log *zap.Logger) *grpcx.Server {
	// 读取配置
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"` // etcd地址
		Name     string `yaml:"name"`     // 服务名
		Port     int    `yaml:"port"`     // 服务端口
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}

	// 创建 grpc server
	grpcServer := grpc.NewServer()
	// 注册服务
	interactiveServer.Register(grpcServer)

	return &grpcx.Server{
		Server:   grpcServer,
		EtcdAddr: cfg.EtcdAddr,
		Name:     cfg.Name,
		Port:     cfg.Port,
		Log:      log,
	}
}
