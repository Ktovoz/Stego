package engine

import (
	"crypto/sha256"
	"encoding/binary"
)

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func scatterParams(password string, n int, context []byte) (a, b int) {
	if n <= 0 {
		return 1, 0
	}
	h := sha256.Sum256(append(append([]byte(password), '|'), context...))
	x := binary.LittleEndian.Uint64(h[0:8])
	a = int((x | 1) % uint64(n))
	if a == 0 {
		a = 1
	}
	for gcd(a, n) != 1 {
		a = (a + 2) % n
		if a == 0 {
			a = 1
		}
	}
	b = int(binary.LittleEndian.Uint64(h[8:16]) % uint64(n))
	return a, b
}

func scatterSlotIndex(k, n, a, b int) int {
	return (a*k + b) % n
}
