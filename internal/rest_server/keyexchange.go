package rest_server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"github.com/globe-and-citizen/layer8-genesis-repo/api/rest"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
	"github.com/go-playground/validator/v10"
)

// pubs stores the public keys of the clients
// TODO: move this to a database
var pubs map[string]struct {
	X, Y *big.Int
}

func init() {
	pubs = make(map[string]struct {
		X, Y *big.Int
	})
}

func (s *RESTServer) keyExchangeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req rest.KeyExchangeRequest
		res rest.KeyExchangeResponse
	)

	// reply is a helper function to send response
	reply := func(status int, statusText string, data []byte) {
		res := &api.Response{
			Status:     status,
			StatusText: statusText,
			Body:       data,
		}
		w.WriteHeader(status)
		resByte, _ := res.ToJSON()
		_, err := w.Write(resByte)
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}

	// parse request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		reply(http.StatusBadRequest, err.Error(), nil)
		return
	}

	// validation
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		reply(http.StatusBadRequest, "validation error", []byte(err.(*validator.InvalidValidationError).Error()))
		return
	}

	// generate nonce
	nonce, err := pkg.GenerateNonce()
	if err != nil {
		log.Printf("Error generating nonce: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil)
		return
	}
	// store public key
	pubX := big.NewInt(0)
	if _, ok := pubX.SetString(req.PublicKeyX, 10); !ok {
		log.Printf("Error parsing public key X: %v", err)
		return
	}
	pubY := big.NewInt(0)
	if _, ok := pubY.SetString(req.PublicKeyY, 10); !ok {
		log.Printf("Error parsing public key Y: %v", err)
		return
	}
	b64nonce := base64.StdEncoding.EncodeToString(nonce)
	pubs[b64nonce] = struct {
		X, Y *big.Int
	}{
		X: pubX,
		Y: pubY,
	}

	// send response
	res.PublicKeyX = s.conf.KeyPair.Pub.X.String()
	res.PublicKeyY = s.conf.KeyPair.Pub.Y.String()
	res.Nonce = b64nonce
	// marshal response
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil)
		return
	}
	reply(http.StatusOK, "ok", resBytes)
}
