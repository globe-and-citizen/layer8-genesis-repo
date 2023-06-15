package grpc_server

import (
	"fmt"
	"net"

	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"github.com/globe-and-citizen/layer8-genesis-repo/internal"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	conf   *config.Config
	server *grpc.Server
}

func NewServer(conf *config.Config) internal.ServerImpl {
	return &GRPCServer{
		conf:   conf,
		server: grpc.NewServer(),
	}
}

// Serve registers the services and starts serving
func (g *GRPCServer) Serve() error {
	lis, err := net.Listen("tcp", ":"+fmt.Sprint(g.conf.GRPCPort))
	if err != nil {
		panic(err)
	}

	// register services
	pb.RegisterKeyExchangeServer(g.server, &KeyExchangeServer{conf: g.conf})
	pb.RegisterDataTransferServiceServer(g.server, &DataTransferServer{conf: g.conf})

	return g.server.Serve(lis)
}

// Shutdown gracefully shuts down the server
func (g *GRPCServer) Shutdown() {
	g.server.GracefulStop()
}
