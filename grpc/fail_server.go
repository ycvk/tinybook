package tinybook

import (
	"context"
	"github.com/cockroachdb/errors/grpc/status"
	"google.golang.org/grpc/codes"
	"math/rand"
)

// FailServer 用来模拟grpc服务端的失败情况
type FailServer struct {
	UserServiceServer
	name string
}

func (s *FailServer) GetUser(ctx context.Context, request *GetUserRequest) (*User, error) {
	name := "fail_" + serverName
	u := &User{
		Id:   1,
		Age:  20,
		Name: &name,
	}
	int31 := rand.Int31()
	if int31%2 == 0 {
		return u, nil
	}
	// 模拟grpc失败 UNAVAILABLE
	return nil, status.Errorf(codes.Unavailable, "模拟服务端失败")
}

//func init() {
//	s2 := uuid.New()
//	serverName = "fail_test_" + s2
//}
