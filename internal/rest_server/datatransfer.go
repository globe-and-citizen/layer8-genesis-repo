package rest_server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"github.com/globe-and-citizen/layer8-genesis-repo/api/rest"
	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
	"github.com/go-playground/validator/v10"
)

func (s *RESTServer) DataTransferHandler(w http.ResponseWriter, r *http.Request) {
	// reply is a helper function to send response
	// reply is a helper function to send response
	reply := func(status int, statusText string, data []byte, headers map[string]string) {
		res := &api.Response{
			Status:     status,
			StatusText: statusText,
			Headers:    headers,
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
	var req rest.DataTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		reply(http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	// validate request
	if err := validator.New().Struct(req); err != nil {
		reply(http.StatusBadRequest, "validation error", []byte(err.(*validator.InvalidValidationError).Error()), nil)
		return
	}

	// get client's public key
	nonce := base64.StdEncoding.EncodeToString(req.Nonce)
	pub, ok := pubs[nonce]
	if !ok { // if nonce not found
		reply(http.StatusUnauthorized, "unauthorized", nil, nil)
		return
	}

	// generate symmetric key
	key, err := pkg.GenerateSharedSecret(s.conf.KeyPair.Pri.D, pub.X, pub.Y)
	if err != nil {
		log.Printf("Error generating symmetric key: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil, nil)
		return
	}

	// decrypt data
	data, err := pkg.DecryptData(key, req.Nonce, req.Data)
	if err != nil {
		log.Printf("Error decrypting data: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil, nil)
	}

	// send data to internet
	reqData, err := api.FromJSONRequest(data)
	if err != nil {
		reply(http.StatusBadRequest, err.Error(), nil, nil)
		return
	}
	httpReq, err := http.NewRequest(reqData.Method, reqData.Url, io.NopCloser(bytes.NewReader(reqData.Body)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil, nil)
		return
	}
	for k, v := range reqData.Headers {
		httpReq.Header.Set(k, v)
	}
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil, nil)
		return
	}
	defer httpRes.Body.Close()

	// read response
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil, nil)
		return
	}
	headers := make(map[string]string)
	for k, v := range httpRes.Header {
		headers[k] = v[0]
	}

	// encrypt response
	encData, err := pkg.EncryptData(key, req.Nonce, body)
	if err != nil {
		log.Printf("Error encrypting response: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil, nil)
		return
	}

	// send response
	var resData rest.DataTransferResponse
	resData.Data = encData
	resByte, err := json.Marshal(&resData)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		reply(http.StatusInternalServerError, "internal server error", nil, nil)
		return
	}
	reply(httpRes.StatusCode, httpRes.Status, resByte, headers)
}
