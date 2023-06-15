package http1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"github.com/globe-and-citizen/layer8-genesis-repo/api/rest"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

func (c *Client) transfer(key string, nonce []byte, req []byte) (*api.Response, error) {
	// encrypt request
	encReq, err := pkg.EncryptData(key, nonce, req)
	if err != nil {
		return nil, err
	}

	// send request
	req, err = json.Marshal(&rest.DataTransferRequest{
		Nonce: nonce,
		Data:  encReq,
	})
	if err != nil {
		return nil, err
	}
	res, err := http.Post(c.baseURL+"/api/v1/data", "application/json", bytes.NewReader(req))
	if err != nil {
		return nil, err
	}

	// parse response
	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp, err := api.FromJSONResponse(byteData)
	if err != nil {
		return nil, err
	}
	var respBody rest.DataTransferResponse
	if err := json.Unmarshal(resp.Body, &respBody); err != nil {
		return nil, err
	}

	// decrypt response
	decRes, err := pkg.DecryptData(key, nonce, respBody.Data)
	if err != nil {
		return nil, err
	}
	resp.Body = decRes
	return resp, nil
}
