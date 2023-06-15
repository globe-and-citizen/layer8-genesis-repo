package pkg

import (
	"crypto/aes"
	"crypto/cipher"
)

// EncryptData encrypts the data using the symmetric key and nonce
func EncryptData(key string, nonce, data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nil, nonce, data, nil), nil
}

// DecryptData decrypts the data using the symmetric key and nonce
func DecryptData(key string, nonce, data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}
