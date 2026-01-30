package crypto

import (
	"math/rand"
	"testing"
)

func TestRSRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	msg := make([]byte, RSK)
	rng.Read(msg)

	cw := rsEncode(msg, RSNSym)
	decoded, err := rsDecode(cw, RSK, RSNSym)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if string(decoded) != string(msg) {
		t.Fatalf("decoded mismatch")
	}
}

func TestECCWrapUnwrapRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	data := make([]byte, 5000)
	rng.Read(data)

	wrapped, err := ECCWrapRS(data)
	if err != nil {
		t.Fatalf("wrap failed: %v", err)
	}

	unwrapped, err := ECCUnwrapRS(wrapped)
	if err != nil {
		t.Fatalf("unwrap failed: %v", err)
	}
	if string(unwrapped) != string(data) {
		t.Fatalf("unwrap mismatch")
	}
}
