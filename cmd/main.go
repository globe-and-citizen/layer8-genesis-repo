package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"github.com/globe-and-citizen/layer8-genesis-repo/internal/grpc_server"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

var port = flag.Int("port", 50055, "The server port. Default: 50055")

func main() {
	flag.Parse()

	// generate key pair for this server instance
	pair, err := pkg.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	// create config
	conf := &config.Config{
		Port:    *port,
		KeyPair: pair,
	}

	// create server listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		panic(err)
	}
	log.Printf("Server listening at %v", lis.Addr())

	// start server
	server := grpc_server.NewServer(conf)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
