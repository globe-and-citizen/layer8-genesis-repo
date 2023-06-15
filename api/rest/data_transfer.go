package rest

type DataTransferRequest struct {
	Nonce []byte `json:"nonce" validate:"required"`
	Data  []byte `json:"data" validate:"required"`
}

type DataTransferResponse struct {
	Data []byte `json:"data" validate:"required"`
}
