package pkg_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := pkg.GenerateKeyPair()
	assert.NoError(t, err)

	pub := kp.Pub
	if pub.X.Cmp(big.NewInt(0)) == 0 || pub.Y.Cmp(big.NewInt(0)) == 0 {
		t.Errorf("Public key is zero")
	}

	assert.NotEqual(t, kp.Pri, big.NewInt(0))
	assert.NotEqual(t, kp.Pub, big.NewInt(0))
}

func TestGenerateNonce(t *testing.T) {
	nonces := make(map[string]bool)

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("generate_nonce_%d", i), func(t *testing.T) {
			nonce, err := pkg.GenerateNonce()
			assert.NoError(t, err)
			if nonces[string(nonce)] {
				t.Errorf("Nonce already exists")
			}
			nonces[string(nonce)] = true
		})
	}
}

func TestGenerateSharedSecret(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("generate_symmetric_key_%d", i), func(t *testing.T) {
			kp1, err := pkg.GenerateKeyPair()
			assert.NoError(t, err)
			kp2, err := pkg.GenerateKeyPair()
			assert.NoError(t, err)

			sec1, err := pkg.GenerateSharedSecret(kp1.Pri.D, kp2.Pub.X, kp2.Pub.Y)
			assert.NoError(t, err)
			sec2, err := pkg.GenerateSharedSecret(kp2.Pri.D, kp1.Pub.X, kp1.Pub.Y)
			assert.NoError(t, err)

			assert.Equal(t, sec1, sec2)
		})
	}
}
