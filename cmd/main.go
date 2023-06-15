package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"github.com/globe-and-citizen/layer8-genesis-repo/internal"
	"github.com/globe-and-citizen/layer8-genesis-repo/internal/grpc_server"
	"github.com/globe-and-citizen/layer8-genesis-repo/internal/rest_server"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

var (
	restPtr = flag.Int("rest-port", 0, "Run REST server on this port")
	grpcPtr = flag.Int("grpc-port", 0, "Run gRPC server on this port")
)

func main() {
	flag.Parse()

	// generate key pair for this server instance
	pair, err := pkg.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	// create config
	conf := &config.Config{
		RESTPort: *restPtr,
		GRPCPort: *grpcPtr,
		KeyPair:  pair,
	}
	// at least one server must be started
	if conf.RESTPort == 0 && conf.GRPCPort == 0 {
		panic("At least one server must be started. Specify '--rest-port' or '--grpc-port' or both")
	}

	// start servers
	var servers []internal.ServerImpl
	// REST
	if conf.RESTPort != 0 {
		log.Printf("REST server is running on port %d", conf.RESTPort)
		server := rest_server.NewServer(conf)
		servers = append(servers, server)
		go func() {
			if err := server.Serve(); err != nil {
				panic(err)
			}
		}()
	}

	// gRPC
	if conf.GRPCPort != 0 {
		log.Printf("GRPC server is running on port %d", conf.GRPCPort)
		server := grpc_server.NewServer(conf)
		servers = append(servers, server)
		go func() {
			if err := server.Serve(); err != nil {
				panic(err)
			}
		}()
	}

	// wait for interrupt signal to shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("Shutting all servers...")
	for _, s := range servers {
		s.Shutdown()
	}
	log.Println("All servers are shut down")
}
