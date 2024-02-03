package grpcx

import (
	"context"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
	"tinybook/tinybook/pkg/netx"
)

type Server struct {
	*grpc.Server
	client     *etcdv3.Client // etcd客户端
	EtcdAddr   string         // etcd地址
	Name       string         // 服务名
	Port       int            // 服务端口
	CancelFunc context.CancelFunc
	Log        *zap.Logger
}

func (s *Server) Serve() error {
	addr := ":" + strconv.Itoa(s.Port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	err = s.register()
	if err != nil {
		return err
	}
	return s.Server.Serve(listen)
}

func (s *Server) register() error {
	cli, err := etcdv3.NewFromURL(s.EtcdAddr)
	if err != nil {
		return err
	}
	s.client = cli

	// 创建一个服务管理器
	manager, err := endpoints.NewManager(s.client, "service/"+s.Name)
	if err != nil {
		return err
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	var ttl int64 = 3
	// 设置服务地址的ttl
	grant, err2 := s.client.Grant(timeout, ttl)
	if err2 != nil {
		return err2
	}

	addr := netx.GetOutboundIP() + ":" + strconv.Itoa(s.Port) // 服务地址
	//addr := "127.0.0.1:" + strconv.Itoa(s.Port) // 服务地址
	key := "service/" + s.Name + "/" + addr // 服务地址对应的键
	// 往etcd中添加一个服务地址 服务名为以上的key
	err = manager.AddEndpoint(
		timeout, key,
		endpoints.Endpoint{Addr: addr},
		etcdv3.WithLease(grant.ID),
	)
	if err != nil {
		return err
	}

	// 服务地址的心跳
	withCancel, c := context.WithCancel(context.Background())
	s.CancelFunc = c
	ch, err3 := s.client.KeepAlive(withCancel, grant.ID)
	if err3 != nil {
		return err3
	}
	go func() {
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					s.Log.Info("keepalive channel closed")
					return
				}
			}
		}
	}()
	return nil
}

func (s *Server) Close() error {
	if s.CancelFunc != nil {
		s.CancelFunc()
	}
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			return err
		}
	}
	s.GracefulStop()
	return nil
}
