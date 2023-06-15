package http2

import (
	"context"
	"math/big"
	"time"

	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

// exchangeKey exchanges a public key with the server
// and returns the client's symmetric key and nonce
func (c *connection) exchangeKey() (string, []byte, error) {
	client := pb.NewKeyExchangeClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// generate key pair
	pair, err := pkg.GenerateKeyPair()
	if err != nil {
		return "", nil, err
	}

	// exchange key
	resp, err := client.ExchangeKey(ctx, &pb.KeyExchangeRequest{
		PublicKeyX: pair.Pub.X.Int64(),
		PublicKeyY: pair.Pub.Y.Int64(),
	})
	if err != nil {
		return "", nil, err
	}

	// generate symmetric key
	symm, err := pkg.GenerateSharedSecret(
		pair.Pri.D, big.NewInt(resp.PublicKeyX), big.NewInt(resp.PublicKeyY))
	if err != nil {
		return "", nil, err
	}
	return symm, resp.Nonce, nil
}
