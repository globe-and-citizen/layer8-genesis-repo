package rest

type KeyExchangeRequest struct {
	PublicKeyX string `json:"public_key_x" validate:"required"`
	PublicKeyY string `json:"public_key_y" validate:"required"`
}

type KeyExchangeResponse struct {
	PublicKeyX string `json:"public_key_x" validate:"required"`
	PublicKeyY string `json:"public_key_y" validate:"required"`
	Nonce      string `json:"nonce" validate:"required"`
}
