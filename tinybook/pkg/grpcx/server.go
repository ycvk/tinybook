package grpcx

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	Address string
}

func (s *Server) Serve() error {
	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	return s.Server.Serve(listen)
}
