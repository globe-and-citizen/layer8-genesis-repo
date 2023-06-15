package http1

import (
	"fmt"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"github.com/globe-and-citizen/layer8-genesis-repo/web/protocols"
)

type Client struct {
	baseURL string
}

// NewClient creates a new instance of the http1.Client
//
// Arguments:
//
//	host (string): The host of the remote layer8 REST server
//	port (string): The port of the remote layer8 REST server
func NewClient(protocol, host, port string) protocols.ClientImpl {
	return &Client{
		baseURL: fmt.Sprintf("%s://%s:%s", protocol, host, port),
	}
}

func (c *Client) Do(req *api.Request) *api.Response {
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
