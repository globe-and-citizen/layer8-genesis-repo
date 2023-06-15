package grpc_server

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	pb "github.com/globe-and-citizen/layer8-genesis-repo/api/grpc"
	"github.com/globe-and-citizen/layer8-genesis-repo/config"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

type DataTransferServer struct {
	pb.UnimplementedDataTransferServiceServer
	conf *config.Config
}

func (s *DataTransferServer) Transfer(ctx context.Context, in *pb.DataTransferRequest) (*pb.DataTransferResponse, error) {
	log.Printf("Received: %v", in)

	newErr := func(status int, statusText string) (*pb.DataTransferResponse, error) {
		res := &api.Response{
			Status:     status,
			StatusText: statusText,
		}
		data, _ := res.ToJSON()
		return &pb.DataTransferResponse{
			Data: data,
		}, nil
	}

	// get client's public key
	pub, ok := pubs[string(in.Nonce)]
	if !ok { // if nonce not found
		return newErr(401, "Unauthorized")
	}

	// generate symmetric key
	key, err := pkg.GenerateSharedSecret(s.conf.KeyPair.Pri.D, pub.X, pub.Y)
	if err != nil {
		log.Printf("Error generating symmetric key: %v", err)
		return newErr(500, "Internal Server Error")
	}

	// decrypt data
	data, err := pkg.DecryptData(key, in.Nonce, in.Data)
	if err != nil {
		log.Printf("Error decrypting data: %v", err)
		return newErr(500, "Internal Server Error")
	}

	// send data to internet
	req, err := api.FromJSONRequest(data)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return newErr(400, err.Error())
	}
	httpReq, err := http.NewRequest(req.Method, req.Url, io.NopCloser(bytes.NewReader(req.Body)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return newErr(400, err.Error())
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return newErr(500, "Internal Server Error")
	}
	defer httpRes.Body.Close()

	// read response
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return newErr(500, "Internal Server Error")
	}

	// create response object
	resObj := &api.Response{
		Status:     httpRes.StatusCode,
		StatusText: httpRes.Status,
		Body:       body,
	}
	resObj.Headers = make(map[string]string)
	for k, v := range httpRes.Header {
		resObj.Headers[k] = v[0]
	}

	// convert response to json
	body, err = resObj.ToJSON()
	if err != nil {
		log.Printf("Error converting response to json: %v", err)
		return newErr(500, "Internal Server Error")
	}

	// encrypt response
	encData, err := pkg.EncryptData(key, in.Nonce, body)
	if err != nil {
		log.Printf("Error encrypting response: %v", err)
		return newErr(500, "Internal Server Error")
	}

	return &pb.DataTransferResponse{
		Data: encData,
	}, nil
}
