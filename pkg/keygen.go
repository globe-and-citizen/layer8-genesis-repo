package pkg

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/hkdf"
)

type KeyPair struct {
	Pri *ecdsa.PrivateKey
	Pub *ecdsa.PublicKey
}

// GenerateKeyPair generates a new key pair using the elliptic curve
func GenerateKeyPair() (*KeyPair, error) {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		Pri: pri,
		Pub: &pri.PublicKey,
	}, nil
}

// GenerateNonce generates a 12 bytes nonce
func GenerateNonce() ([]byte, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateSharedSecret generates a shared secret using the private and public key
func GenerateSharedSecret(priD, pubX, pubY *big.Int) (string, error) {
	pub := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     pubX,
		Y:     pubY,
	}

	// 32 bytes (AES-256) key
	shared, _ := pub.Curve.ScalarMult(pubX, pubY, priD.Bytes())
	sym := hkdf.New(sha256.New, shared.Bytes(), nil, nil)
	key := make([]byte, 32)
	if _, err := sym.Read(key); err != nil {
		return "", err
	}
	return string(key), nil
}
