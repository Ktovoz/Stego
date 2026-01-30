package crypto

import (
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
)

func PBKDF2Compat(password string, salt []byte, iterations int, keyLen int) []byte {
	return pbkdf2Key([]byte(password), salt, iterations, keyLen)
}

func pbkdf2Key(password []byte, salt []byte, iterations int, keyLen int) []byte {
	return pbkdf2.Key(password, salt, iterations, keyLen, sha1.New)
}
