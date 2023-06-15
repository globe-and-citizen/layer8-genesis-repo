package protocols

import "github.com/globe-and-citizen/layer8-genesis-repo/api"

type ClientImpl interface {
	// Do sends a request through an encrypted channel over HTTP/2 and returns a response
	Do(req *api.Request) *api.Response
}
