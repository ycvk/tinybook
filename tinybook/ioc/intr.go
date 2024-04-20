package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "tinybook/tinybook/api/proto/gen/intr/v1"
	"tinybook/tinybook/interactive/service"
	client2 "tinybook/tinybook/internal/client"
)

// InitIntrClientV1 初始化交互服务客户端 用于远程服务的调用
func InitIntrClientV1(client3 *etcdv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	// 创建一个服务管理器
	builder, err := resolver.NewBuilder(client3)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{
		grpc.WithResolvers(builder), // 注册etcd服务 用于服务发现
	}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // 无证书
	conn, err := grpc.Dial(cfg.Addr, opts...)                                     // 连接远程服务
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(conn)
	return remote
}

// InitIntrClient 初始化交互服务客户端 用于本地和远程服务的联合调用
func InitIntrClient(service service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string `yaml:"addr"`
		Threshold int32  `yaml:"threshold"` // 阈值 用于判断是调用本地服务还是远程服务
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
