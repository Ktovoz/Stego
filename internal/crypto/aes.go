package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type AESGCMConfig struct {
	KeyLength  int
	SaltLength int
	NonceLen   int
	TagLen     int
	Iterations int
}

func DefaultAESGCMConfig() AESGCMConfig {
	return AESGCMConfig{
		KeyLength:  32,
		SaltLength: 16,
		NonceLen:   12,
		TagLen:     16,
		Iterations: 50000,
	}
}

func EncryptAESGCM(key []byte, nonce []byte, plaintext []byte) (ciphertext []byte, tag []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, nil, errors.New("invalid nonce length")
	}
	combined := gcm.Seal(nil, nonce, plaintext, nil)
	if len(combined) < gcm.Overhead() {
		return nil, nil, errors.New("ciphertext too short")
	}
	tag = combined[len(combined)-gcm.Overhead():]
	ciphertext = combined[:len(combined)-gcm.Overhead()]
	return ciphertext, tag, nil
}

func DecryptAESGCM(key []byte, nonce []byte, ciphertext []byte, tag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, errors.New("invalid nonce length")
	}
	combined := append(append([]byte{}, ciphertext...), tag...)
	return gcm.Open(nil, nonce, combined, nil)
}

func RandomBytes(n int) ([]byte, error) {
	if n <= 0 {
		return []byte{}, nil
	}
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	return b, err
}
