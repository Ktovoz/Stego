package engine

import (
	"encoding/binary"
	"errors"
)

func extractBytes2bitAtSlot(rgb []byte, startSlot int, byteLen int) []byte {
	if byteLen <= 0 || startSlot < 0 || startSlot >= len(rgb) {
		return nil
	}
	slotsNeeded := byteLen * 4
	available := len(rgb) - startSlot
	if slotsNeeded > available {
		slotsNeeded = available
		byteLen = slotsNeeded / 4
	}
	if byteLen <= 0 {
		return nil
	}
	out := make([]byte, byteLen)
	slot := 0
	for i := 0; i < byteLen; i++ {
		var b byte
		for shift := 6; shift >= 0; shift -= 2 {
			two := rgb[startSlot+slot] & 0x3
			b |= two << uint(shift)
			slot++
		}
		out[i] = b
	}
	return out
}

func extractScatteredBytes2bit(rgb []byte, startSlot int, byteLen int, password string) []byte {
	if byteLen <= 0 || startSlot < 0 || startSlot >= len(rgb) {
		return nil
	}
	available := len(rgb) - startSlot
	if available <= 0 {
		return nil
	}
	if byteLen*4 > available {
		byteLen = available / 4
	}
	if byteLen <= 0 {
		return nil
	}
	a, b := scatterParams(password, available, []byte("scatter_body_v1"))
	out := make([]byte, byteLen)
	k := 0
	for i := 0; i < byteLen; i++ {
		var bt byte
		for shift := 6; shift >= 0; shift -= 2 {
			idx := startSlot + scatterSlotIndex(k, available, a, b)
			two := rgb[idx] & 0x3
			bt |= two << uint(shift)
			k++
		}
		out[i] = bt
	}
	return out
}

func (e *Engine) Extract(rgb []byte, width, height int, password string) ([]byte, bool, bool, []byte, error) {
	if len(rgb) != width*height*3 {
		return nil, false, false, nil, errors.New("invalid rgb buffer size")
	}
	headerBytes := extractBytes2bitAtSlot(rgb, 0, HeaderLength)
	if len(headerBytes) != HeaderLength {
		return nil, false, false, nil, errors.New("invalid header")
	}
	rawLen := binary.LittleEndian.Uint32(headerBytes)
	integrityEnabled := (rawLen & IntegrityFlag) != 0
	scatterEnabled := (rawLen & ScatterFlag) != 0
	dataLen := int(rawLen & ^uint32(IntegrityFlag|ScatterFlag))
	if !integrityEnabled {
		dataLen = int(rawLen)
	}

	maxSize := e.CalculateMaxCapacity(width, height, true)
	maxSize = maxSize - HeaderLength - CRCLength + 32
	if integrityEnabled {
		maxSize -= IntegrityHashLen
	}
	if dataLen <= 0 || dataLen > maxSize {
		return nil, integrityEnabled, scatterEnabled, nil, errors.New("invalid data length")
	}

	var integrityBytes []byte
	if integrityEnabled {
		integritySlotStart := (HeaderLength * 8) / 2
		integrityBytes = extractBytes2bitAtSlot(rgb, integritySlotStart, IntegrityHashLen)
		if len(integrityBytes) != IntegrityHashLen {
			return nil, integrityEnabled, scatterEnabled, nil, errors.New("invalid integrity")
		}
	}

	fixedLen := HeaderLength
	if integrityEnabled {
		fixedLen += IntegrityHashLen
	}
	startSlot := (fixedLen * 8) / 2

	var extractedData []byte
	var crcBytes []byte
	if scatterEnabled {
		if password == "" {
			return nil, integrityEnabled, scatterEnabled, integrityBytes, errors.New("password required for scattered data")
		}
		body := extractScatteredBytes2bit(rgb, startSlot, dataLen+CRCLength, password)
		if len(body) != dataLen+CRCLength {
			return nil, integrityEnabled, scatterEnabled, integrityBytes, errors.New("invalid scattered payload length")
		}
		extractedData = body[:dataLen]
		crcBytes = body[dataLen:]
	} else {
		body := extractBytes2bitAtSlot(rgb, startSlot, dataLen+CRCLength)
		if len(body) != dataLen+CRCLength {
			return nil, integrityEnabled, scatterEnabled, integrityBytes, errors.New("invalid payload length")
		}
		extractedData = body[:dataLen]
		crcBytes = body[dataLen:]
	}
	if !verifyCRC32(extractedData, crcBytes) {
		return nil, integrityEnabled, scatterEnabled, integrityBytes, errors.New("crc32 verify failed")
	}

	if integrityEnabled && len(integrityBytes) == IntegrityHashLen {
		actual, err := embeddedPixelHash(rgb, width, height)
		if err == nil {
			_ = equalBytes(actual, integrityBytes)
		}
	}

	return extractedData, integrityEnabled, scatterEnabled, integrityBytes, nil
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
