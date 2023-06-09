package http2

import (
	"fmt"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"github.com/globe-and-citizen/layer8-genesis-repo/web/protocols"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	host string
	port string
	conn *grpc.ClientConn
}

// NewClient creates a new instance of the http2.Client
//
// Arguments:
//
//	host (string): The host of the remote layer8 gRPC server
//	port (string): The port of the remote layer8 gRPC server
func NewClient(host, port string) protocols.ClientImpl {
	return &Client{
		host: host,
		port: port,
	}
}

func (c *Client) Do(req *api.Request) *api.Response {
	// connect to the server
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", c.host, c.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &api.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	defer conn.Close()
	c.conn = conn

	// exchange key
	key, nonce, err := c.exchangeKey()
	if err != nil {
		return &api.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}

	// send request
	jsonReq, err := req.ToJSON()
	if err != nil {
		return &api.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	res, err := c.transfer(key, nonce, jsonReq)
	if err != nil {
		return &api.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	return res
}
