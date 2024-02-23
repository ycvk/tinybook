package tinybook

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
	_ "tinybook/tinybook/pkg/grpcx/balancer/wrr"
)

type EtcdBalancerTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (s *EtcdBalancerTestSuite) SetupSuite() {
	var err error
	s.client, err = etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:32379"},
	})
	s.Require().NoError(err)
}

func (s *EtcdBalancerTestSuite) TestServer() {
	port := GenFreePort()
	listen, err2 := net.Listen("tcp", fmt.Sprintf(":%d", port))
	s.Require().NoError(err2)
	t := s.T()
	// 创建一个服务管理器
	manager, err := endpoints.NewManager(s.client, "service/user")
	require.NoError(t, err)
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	var ttl int64 = 3
	// 设置服务地址的ttl
	grant, err2 := s.client.Grant(timeout, ttl)
	require.NoError(t, err2)

	addr := "127.0.0.1" + fmt.Sprintf(":%d", port) // 服务地址
	key := "service/user/" + addr                  // 服务地址对应的键
	// 往etcd中添加一个服务地址 服务名为service/user
	m := make(map[string]any)
	m["weight"] = 1
	err = manager.AddEndpoint(
		timeout, key,
		endpoints.Endpoint{Addr: addr, Metadata: m},
		etcdv3.WithLease(grant.ID),
	)
	require.NoError(t, err)

	withCancel, c := context.WithCancel(context.Background())
	go func() {
		alive, err3 := s.client.KeepAlive(withCancel, grant.ID)
		require.NoError(t, err3)
		for {
			select {
			case _, ok := <-alive:
				if !ok {
					t.Log("keepalive channel closed")
					return
				}
			}

		}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{})
	err = server.Serve(listen)
	if err != nil {
		c() // 取消etcd的keepalive
		return
	}
	c() // 取消etcd的keepalive

	err2 = manager.DeleteEndpoint(context.Background(), key) // 删除服务地址
	require.NoError(t, err2)
	server.GracefulStop()
}

// GenFreePort 获取一个空闲的端口;端口避免写死,因为要启动多个实例,测试负载均衡
func GenFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listen.Close()
	return listen.Addr().(*net.TCPAddr).Port
}

func (s *EtcdBalancerTestSuite) TestClient() {
	t := s.T()
	builder, err := resolver.NewBuilder(s.client)
	require.NoError(t, err)

	conn, err := grpc.Dial("etcd:///service/user",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
  "loadBalancingConfig": [
    {
      "custom_weighted_round_robin": {}
    }
  ],
  "methodConfig": [
    {
      "name": [
        {
          "service": "UserService"
        }
      ],
      "retryPolicy": {
        "maxAttempts": 3,
        "initialBackoff": "0.01s",
        "maxBackoff": "0.1s",
        "backoffMultiplier": 1.3,
        "retryableStatusCodes": [
          "UNAVAILABLE"
        ]
      }
    }
  ]
}`),
		grpc.WithResolvers(builder), // 注册etcd的resolver
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewUserServiceClient(conn)
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	m := make(map[string]int)
	for i := 0; i < 5000; i++ {
		user, err := client.GetUser(timeout, &GetUserRequest{
			Id: 1,
		})
		if err != nil {
			//t.Logf("err:%v", err)
			m[err.Error()]++
		} else {
			//t.Log(user)
		}
		if user != nil {
			m[*user.Name]++
		}
	}
	for k, v := range m {
		t.Log(k, v)
	}
}

func (s *EtcdBalancerTestSuite) TestStartMultiServers() {
	go s.TestFailServer()
	s.TestServer()
}

func TestEtcdBalancer(t *testing.T) {
	suite.Run(t, new(EtcdBalancerTestSuite))
}

func (s *EtcdBalancerTestSuite) TestFailServer() {
	port := GenFreePort()
	listen, err2 := net.Listen("tcp", fmt.Sprintf(":%d", port))
	s.Require().NoError(err2)
	t := s.T()
	// 创建一个服务管理器
	manager, err := endpoints.NewManager(s.client, "service/user")
	require.NoError(t, err)
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	var ttl int64 = 3
	// 设置服务地址的ttl
	grant, err2 := s.client.Grant(timeout, ttl)
	require.NoError(t, err2)

	addr := "127.0.0.1" + fmt.Sprintf(":%d", port) // 服务地址
	key := "service/user/" + addr                  // 服务地址对应的键
	// 往etcd中添加一个服务地址 服务名为service/user
	m := make(map[string]any)
	m["weight"] = 1
	err = manager.AddEndpoint(
		timeout, key,
		endpoints.Endpoint{Addr: addr, Metadata: m},
		etcdv3.WithLease(grant.ID),
	)
	require.NoError(t, err)

	withCancel, c := context.WithCancel(context.Background())
	go func() {
		alive, err3 := s.client.KeepAlive(withCancel, grant.ID)
		require.NoError(t, err3)
		for {
			select {
			case _, ok := <-alive:
				if !ok {
					t.Log("keepalive channel closed")
					return
				}
			}

		}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &FailServer{})
	err = server.Serve(listen)
	if err != nil {
		c() // 取消etcd的keepalive
		return
	}
	c() // 取消etcd的keepalive

	err2 = manager.DeleteEndpoint(context.Background(), key) // 删除服务地址
	require.NoError(t, err2)
	server.GracefulStop()
}
