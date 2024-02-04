package tinybook

import (
	"context"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	uuid "github.com/lithammer/shortuuid/v4"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"net"
	"testing"
)

const (
	EtcdAddr   = "127.0.0.1"   // etcd地址
	EtcdPort   = "8500"        // etcd端口
	ServerAddr = "127.0.0.1"   // grpc服务地址
	ServerPort = 18080         // grpc服务端口
	ServerName = "grpcService" // grpc服务名
	ServerTag  = "grpc"        // grpc服务在consul中的标签
)

type ConsulTestSuite struct {
	suite.Suite
	client *consulapi.Client // consul客户端
}

func (s *ConsulTestSuite) SetupSuite() {
	config := consulapi.DefaultConfig()
	config.Address = EtcdAddr + ":" + EtcdPort
	// 创建consul客户端
	client, err := consulapi.NewClient(config)
	s.NoError(err)
	s.client = client
}

func (s *ConsulTestSuite) TestConsulServer() {
	listen, err2 := net.Listen("tcp", ":"+fmt.Sprintf("%d", ServerPort))
	s.Require().NoError(err2)
	// 注册grpc服务和grpc健康检查服务
	server := grpc.NewServer()
	healthCheckServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthCheckServer)
	RegisterUserServiceServer(server, &Server{})

	// 创建consul注册对象
	registration := &consulapi.AgentServiceRegistration{
		Address: ServerAddr,
		Port:    ServerPort,
		ID:      uuid.New(), // 服务id 保证唯一
		Name:    ServerName,
		Tags:    []string{ServerTag},
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", ServerAddr, ServerPort),
			Interval:                       "5s",  // 健康检查间隔
			Timeout:                        "2s",  // 健康检查超时
			DeregisterCriticalServiceAfter: "15s", // 故障检查失败15s后 consul自动将注册服务删除
		},
	}

	// 注册服务到consul
	err := s.client.Agent().ServiceRegister(registration)
	s.NoError(err)

	err2 = server.Serve(listen)
	s.Require().NoError(err2)
}

func (s *ConsulTestSuite) TestConsulClient() {
	target := fmt.Sprintf("consul://%s:%s/%s?wait=10s&tag=%s", EtcdAddr, EtcdPort, ServerName, ServerTag)
	dial, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 无证书
		//grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 轮询
	)
	s.Require().NoError(err)
	defer dial.Close()
	// 创建grpc客户端
	client := NewUserServiceClient(dial)
	// 调用grpc服务
	user, err := client.GetUser(context.Background(), &GetUserRequest{Id: 1})
	s.Require().NoError(err)
	s.T().Log(user)
}

func TestConsul(t *testing.T) {
	suite.Run(t, new(ConsulTestSuite))
}
