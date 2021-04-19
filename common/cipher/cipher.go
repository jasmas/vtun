package cipher

import (
	"crypto/sha256"

	"golang.org/x/crypto/chacha20poly1305"
)

var nonce = make([]byte, chacha20poly1305.NonceSizeX)

var hashKey = []byte("8pUsXuZw4z6B9EhGdKgNjQnjmVsYv2x5")

func GenerateKey(key string) {
	sha := sha256.Sum256([]byte(key))
	buff := make([]byte, 32)
	copy(sha[:32], buff[:32])
	hashKey = buff
}

func Encrypt(data *[]byte) {
	aead, _ := chacha20poly1305.NewX(hashKey)
	ciphertext := aead.Seal(nil, nonce, *data, nil)
	data = &ciphertext
}

func Decrypt(data *[]byte) {
	aead, _ := chacha20poly1305.NewX(hashKey)
	plaintext, _ := aead.Open(nil, nonce, *data, nil)
	data = &plaintext
}
