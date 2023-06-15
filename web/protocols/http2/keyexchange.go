package http2

import (
	"context"
	"encoding/base64"
	"math/big"
	"time"

	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

// exchangeKey exchanges a public key with the server
// and returns the client's symmetric key and nonce
func (c *Client) exchangeKey() (string, []byte, error) {
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
		PublicKeyX: pair.Pub.X.String(),
		PublicKeyY: pair.Pub.Y.String(),
	})
	if err != nil {
		return "", nil, err
	}

	// generate symmetric key
	nonce, err := base64.StdEncoding.DecodeString(resp.Nonce)
	if err != nil {
		return "", nil, err
	}
	pubX := big.NewInt(0)
	if _, ok := pubX.SetString(resp.PublicKeyX, 10); !ok {
		return "", nil, err
	}
	pubY := big.NewInt(0)
	if _, ok := pubY.SetString(resp.PublicKeyY, 10); !ok {
		return "", nil, err
	}
	symm, err := pkg.GenerateSharedSecret(pair.Pri.D, pubX, pubY)
	if err != nil {
		return "", nil, err
	}
	return symm, nonce, nil
}
