package tinybook

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"testing"
	"time"
)

type EtcdTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (s *EtcdTestSuite) SetupSuite() {
	var err error
	s.client, err = etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:32379"},
	})
	s.Require().NoError(err)
}

func (s *EtcdTestSuite) TestServer() {
	listen, err2 := net.Listen("tcp", ":8090")
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

	addr := "127.0.0.1:8090"      // 服务地址
	key := "service/user/" + addr // 服务地址对应的键
	// 往etcd中添加一个服务地址 服务名为service/user
	err = manager.AddEndpoint(
		timeout, key,
		endpoints.Endpoint{Addr: addr},
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

	go func() {
		// 模拟metadata信息变动
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				withTimeout, c := context.WithTimeout(context.Background(), time.Second)
				//err3 := manager.Update(withTimeout, []*endpoints.UpdateWithOpts{
				//	{
				//		Update: endpoints.Update{
				//			Op:  endpoints.Add,
				//			Key: key,
				//			Endpoint: endpoints.Endpoint{
				//				Addr:     addr,
				//				Metadata: time.Now().Second(),
				//			},
				//		},
				//	},
				//})
				err3 := manager.AddEndpoint(withTimeout, key, endpoints.Endpoint{Addr: addr,
					Metadata: time.Now().String()}, etcdv3.WithLease(grant.ID))
				c()
				require.NoError(t, err3)
			}
		}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{})
	server.Serve(listen)

	c()                                                      // 取消etcd的keepalive
	err2 = manager.DeleteEndpoint(context.Background(), key) // 删除服务地址
	require.NoError(t, err2)
	server.GracefulStop()
	//s.client.Close() // 关闭etcd连接
}

func (s *EtcdTestSuite) TestClient() {
	t := s.T()
	builder, err := resolver.NewBuilder(s.client)
	require.NoError(t, err)

	conn, err := grpc.Dial("etcd:///service/user",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "custom_weighted_round_robin"}`),
		grpc.WithResolvers(builder), // 注册etcd的resolver
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewUserServiceClient(conn)
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	user, err := client.GetUser(timeout, &GetUserRequest{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	t.Log(user)
	stream, err := client.GetUserList(context.Background(), &GetUserListRequest{
		Ids: []int64{1, 2, 3},
	})
	if err != nil {
		panic(err)
	}
	for {
		user, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		t.Log(user)
	}
}

func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}
