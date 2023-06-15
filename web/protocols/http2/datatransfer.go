package http2

import (
	"context"
	"time"

	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

// transfer transfers data between the client and the server
// and returns the response
func (c *connection) transfer(key string, nonce, data []byte) ([]byte, error) {
	client := pb.NewDataTransferServiceClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// encrypt data
	encData, err := pkg.EncryptData(key, nonce, data)
	if err != nil {
		return nil, err
	}

	// transfer data
	resp, err := client.Transfer(ctx, &pb.DataTransferRequest{
		Nonce: nonce,
		Data:  encData,
	})
	if err != nil {
		return nil, err
	}

	// decrypt response data
	decData, err := pkg.DecryptData(key, nonce, resp.Data)
	if err != nil {
		return nil, err
	}
	return decData, nil
}
