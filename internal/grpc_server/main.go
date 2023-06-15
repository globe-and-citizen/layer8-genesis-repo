package grpc_server

import (
	"net"

	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	conf *config.Config
}

type GRPCServerImpl interface {
	Serve(lis net.Listener) error
}

func NewServer(conf *config.Config) GRPCServerImpl {
	return &GRPCServer{conf}
}

// Serve registers the services and starts serving
func (g *GRPCServer) Serve(lis net.Listener) error {
	s := grpc.NewServer()

	// register services
	pb.RegisterKeyExchangeServer(s, &KeyExchangeServer{conf: g.conf})
	pb.RegisterDataTransferServiceServer(s, &DataTransferServer{conf: g.conf})

	return s.Serve(lis)
}
