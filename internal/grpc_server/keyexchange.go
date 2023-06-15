package grpc_server

import (
	"context"
	"encoding/base64"
	"log"
	"math/big"

	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

// pubs stores the public keys of the clients
// TODO: move this to a database
var pubs map[string]struct {
	X, Y *big.Int
}

func init() {
	pubs = make(map[string]struct {
		X, Y *big.Int
	})
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
	pubX := big.NewInt(0)
	if _, ok := pubX.SetString(in.PublicKeyX, 10); !ok {
		log.Printf("Error parsing public key X: %v", err)
		return nil, err
	}
	pubY := big.NewInt(0)
	if _, ok := pubY.SetString(in.PublicKeyY, 10); !ok {
		log.Printf("Error parsing public key Y: %v", err)
		return nil, err
	}
	b64nonce := base64.StdEncoding.EncodeToString(nonce)
	pubs[b64nonce] = struct {
		X, Y *big.Int
	}{
		X: pubX,
		Y: pubY,
	}

	// convert server's public key to string
	return &pb.KeyExchangeResponse{
		PublicKeyX: s.conf.KeyPair.Pub.X.String(),
		PublicKeyY: s.conf.KeyPair.Pub.Y.String(),
		Nonce:      b64nonce,
	}, nil
}
