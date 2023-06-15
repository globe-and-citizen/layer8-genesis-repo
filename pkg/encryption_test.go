package pkg_test

import (
	"testing"

	"github.com/globe-and-citizen/layer8-genesis-repo/pkg"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecryptData(t *testing.T) {
	tc := []struct {
		name string
		data []byte
	}{
		{
			name: "test1",
			data: []byte("test1"),
		},
		{
			name: "test2",
			data: []byte("hello world"),
		},
		{
			name: "test3",
			data: []byte("this is a test"),
		},
		{
			name: "test4",
			data: []byte("this is a test with a very long string"),
		},
		{
			name: "test5",
			data: []byte("this is a test with a very long string and some special characters !@#$%^&*()_+"),
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			nonce, err := pkg.GenerateNonce()
			assert.NoError(t, err)

			p1_kp, err := pkg.GenerateKeyPair()
			assert.NoError(t, err)
			p2_kp, err := pkg.GenerateKeyPair()
			assert.NoError(t, err)

			// point 1's shared secret
			p1_ss, err := pkg.GenerateSharedSecret(p1_kp.Pri.D, p2_kp.Pub.X, p2_kp.Pub.Y)
			assert.NoError(t, err)
			// point 2's shared secret
			p2_ss, err := pkg.GenerateSharedSecret(p2_kp.Pri.D, p1_kp.Pub.X, p1_kp.Pub.Y)
			assert.NoError(t, err)

			// data encryption @point 1
			encrypted, err := pkg.EncryptData(p1_ss, nonce, tt.data)
			assert.NoError(t, err)
			assert.NotEqual(t, tt.data, encrypted)

			// data decryption @point 2
			decrypted, err := pkg.DecryptData(p2_ss, nonce, encrypted)
			assert.NoError(t, err)
			assert.Equal(t, tt.data, decrypted)
		})
	}
}
