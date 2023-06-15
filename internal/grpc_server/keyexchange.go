package grpc_server

import (
	"context"
	"log"
	"math/big"

	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

// pubs stores the public keys of the clients
var pubs map[string]struct {
	X, Y *big.Int
}

type KeyExchangeServer struct {
	pb.UnimplementedKeyExchangeServer
	conf *config.Config
}

func (s *KeyExchangeServer) ExchangeKey(ctx context.Context, in *pb.KeyExchangeRequest) (*pb.KeyExchangeResponse, error) {
	log.Printf("Received: %v", in)

	// generate nonce
	nonce, err := pkg.GenerateNonce()
	if err != nil {
		log.Printf("Error generating nonce: %v", err)
		return nil, err
	}

	// store public key
	pubs[string(nonce)] = struct {
		X, Y *big.Int
	}{
		X: big.NewInt(in.PublicKeyX),
		Y: big.NewInt(in.PublicKeyY),
	}

	return &pb.KeyExchangeResponse{
		PublicKeyX: s.conf.KeyPair.Pub.X.Int64(),
		PublicKeyY: s.conf.KeyPair.Pub.Y.Int64(),
		Nonce:      nonce,
	}, nil
}
