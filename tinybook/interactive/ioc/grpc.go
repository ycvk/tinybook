package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "tinybook/tinybook/interactive/grpc"
	"tinybook/tinybook/pkg/grpcx"
)

func InitGrpcServer(interactiveServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	// 读取配置
	type Config struct {
		Addr string `yaml:"addr"`
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
		Server:  grpcServer,
		Address: cfg.Addr,
	}
}
