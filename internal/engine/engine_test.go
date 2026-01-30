package engine

import "testing"

func TestHideExtractRoundTrip(t *testing.T) {
	w, h := 256, 256
	rgb := make([]byte, w*h*3)
	for i := range rgb {
		rgb[i] = byte(i * 31)
	}
	payload := []byte("hello world - stego test")

	eng := New(1024 * 1024)
	out, _, err := eng.Hide(rgb, w, h, payload, "pass", true)
	if err != nil {
		t.Fatalf("hide failed: %v", err)
	}
	got, _, _, _, err := eng.Extract(out, w, h, "pass")
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("payload mismatch")
	}
}

func TestHideExtractRoundTrip_NoScatter(t *testing.T) {
	w, h := 128, 128
	rgb := make([]byte, w*h*3)
	for i := range rgb {
		rgb[i] = byte(i * 17)
	}
	payload := make([]byte, 2048)
	for i := range payload {
		payload[i] = byte(i * 13)
	}

	eng := New(1024 * 1024)
	out, _, err := eng.Hide(rgb, w, h, payload, "", false)
	if err != nil {
		t.Fatalf("hide failed: %v", err)
	}
	got, integrityEnabled, scatterEnabled, _, err := eng.Extract(out, w, h, "")
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}
	if !integrityEnabled {
		t.Fatalf("expected integrity enabled")
	}
	if scatterEnabled {
		t.Fatalf("expected scatter disabled")
	}
	if len(got) != len(payload) {
		t.Fatalf("payload length mismatch")
	}
	for i := range payload {
		if got[i] != payload[i] {
			t.Fatalf("payload mismatch at %d", i)
		}
	}
}
