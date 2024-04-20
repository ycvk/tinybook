package tinybook

import (
	"context"
	uuid "github.com/lithammer/shortuuid/v4"
)

var serverName = "UserService"

type Server struct {
	UserServiceServer
	name string
}

func (s *Server) GetUser(ctx context.Context, request *GetUserRequest) (*User, error) {
	name := serverName
	u := &User{
		Id:   1,
		Age:  20,
		Name: &name,
	}
	return u, nil
}

func (s *Server) GetUserList(request *GetUserListRequest, server UserService_GetUserListServer) error {
	name := "test"

	users := []User{
		{
			Id:   1,
			Age:  20,
			Name: &name,
		},
		{
			Id:   2,
			Age:  24,
			Name: &name,
		},
		{
			Id:   3,
			Age:  25,
			Name: &name,
		},
	}
	for i := range users {
		err := server.Send(&users[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	s2 := uuid.New()
	serverName = "test_" + s2
}
