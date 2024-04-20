package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/interactive/service"
	client2 "tinybook/tinybook/internal/client"
)

func InitIntrClient(service service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string `yaml:"addr"`
		Threshold int32  `yaml:"threshold"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // 无证书
	conn, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(conn)
	local := client2.NewLocalInteractiveServiceAdapter(service)
	interactiveClient := client2.NewInteractiveClient(remote, local, cfg.Threshold)

	// 监听配置变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg := Config{}
		err := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {
			panic(err)
		}
		interactiveClient.SetThreshold(cfg.Threshold)

	})
	return interactiveClient
}
