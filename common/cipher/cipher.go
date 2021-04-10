package cipher

import (
	"crypto/sha256"

	"golang.org/x/crypto/chacha20poly1305"
)

var nonce = make([]byte, chacha20poly1305.NonceSizeX)

func CreateHash(key string) []byte {
	sha := sha256.Sum256([]byte(key))
	ret := make([]byte, 32)
	copy(sha[:32], ret[:32])
	return ret
}

func Encrypt(data *[]byte, key []byte) {
	aead, _ := chacha20poly1305.NewX(key)
	ciphertext := aead.Seal(nil, nonce, *data, nil)
	data = &ciphertext
}

func Decrypt(data *[]byte, key []byte) {
	aead, _ := chacha20poly1305.NewX(key)
	plaintext, _ := aead.Open(nil, nonce, *data, nil)
	data = &plaintext
}
