package http1

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"github.com/globe-and-citizen/layer8-genesis-repo/api/rest"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
)

func (c *Client) exchangeKey() (string, []byte, error) {
	// generate key pair
	pair, err := pkg.GenerateKeyPair()
	if err != nil {
		return "", nil, err
	}

	// exchange key
	req, err := json.Marshal(&rest.KeyExchangeRequest{
		PublicKeyX: pair.Pub.X.String(),
		PublicKeyY: pair.Pub.Y.String(),
	})
	if err != nil {
		return "", nil, err
	}
	res, err := http.Post(c.baseURL+"/api/v1/keys", "application/json", bytes.NewReader(req))
	if err != nil {
		return "", nil, err
	}

	// parse response
	var resp api.Response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", nil, err
	}

	var respBody rest.KeyExchangeResponse
	if err := json.Unmarshal(resp.Body, &respBody); err != nil {
		return "", nil, err
	}

	// generate symmetric key
	pubX := big.NewInt(0)
	if _, ok := pubX.SetString(respBody.PublicKeyX, 10); !ok {
		return "", nil, err
	}
	pubY := big.NewInt(0)
	if _, ok := pubY.SetString(respBody.PublicKeyY, 10); !ok {
		return "", nil, err
	}
	symm, err := pkg.GenerateSharedSecret(pair.Pri.D, pubX, pubY)
	if err != nil {
		return "", nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(respBody.Nonce)
	if err != nil {
		return "", nil, err
	}
	return symm, nonce, nil
}
