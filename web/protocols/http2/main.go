package http2

import (
	"fmt"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	// HTTP2Client is used to make requests to the server over the HTTP2 protocol, with the
	// requests being sent over gRPC
	Client struct {
		host string
		port string
	}

	// connection is used to store the connection to the server
	connection struct {
		conn *grpc.ClientConn
	}
)

type ClientImpl interface {
	// Do sends a request through an encrypted channel over HTTP/2 and returns a response
	Do(req *api.Request) *api.Response
}

// NewClient creates a new instance of the http2.Client
//
// Arguments:
//
//	host (string): The host of the remote layer8 server
//	port (string): The port of the remote layer8 server
func NewClient(host, port string) ClientImpl {
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

	// create a new client
	ct := &connection{conn: conn}

	// exchange key
	key, nonce, err := ct.exchangeKey()
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
	res, err := ct.transfer(key, nonce, jsonReq)
	if err != nil {
		return &api.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	jsonRes, err := api.FromJSONResponse(res)
	if err != nil {
		return &api.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	return jsonRes
}
