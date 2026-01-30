package engine

import (
	"encoding/binary"
	"errors"
)

func embedBytes2bitAtSlot(rgb []byte, startSlot int, data []byte) {
	if startSlot < 0 || startSlot >= len(rgb) || len(data) == 0 {
		return
	}
	slotsNeeded := len(data) * 4
	if startSlot+slotsNeeded > len(rgb) {
		slotsNeeded = len(rgb) - startSlot
	}
	slot := 0
	for i := 0; i < len(data) && slot < slotsNeeded; i++ {
		b := data[i]
		for shift := 6; shift >= 0 && slot < slotsNeeded; shift -= 2 {
			two := (b >> uint(shift)) & 0x3
			idx := startSlot + slot
			rgb[idx] = (rgb[idx] & 0xFC) | two
			slot++
		}
	}
}

func embedScatteredBytes2bit(rgb []byte, startSlot int, data []byte, password string) error {
	if startSlot < 0 {
		startSlot = 0
	}
	available := len(rgb) - startSlot
	if available <= 0 {
		return errors.New("image capacity insufficient: no writable area")
	}
	slotsNeeded := len(data) * 4
	if slotsNeeded > available {
		return errors.New("image capacity insufficient: data exceeds available area")
	}
	a, b := scatterParams(password, available, []byte("scatter_body_v1"))
	k := 0
	for i := 0; i < len(data); i++ {
		bt := data[i]
		for shift := 6; shift >= 0; shift -= 2 {
			two := (bt >> uint(shift)) & 0x3
			idx := startSlot + scatterSlotIndex(k, available, a, b)
			rgb[idx] = (rgb[idx] & 0xFC) | two
			k++
		}
	}
	return nil
}

func (e *Engine) Hide(rgb []byte, width, height int, data []byte, password string, scatter bool) ([]byte, []byte, error) {
	maxCap := e.CalculateMaxCapacity(width, height, false)
	totalBitsNeeded := (HeaderLength + IntegrityHashLen + len(data) + CRCLength) * 8
	if totalBitsNeeded > maxCap*8 {
		return nil, nil, errors.New("image capacity insufficient")
	}

	crcBytes := calculateCRC32(data)
	flags := uint32(IntegrityFlag)
	scatterEnabled := password != "" && scatter
	if scatterEnabled {
		flags |= ScatterFlag
	}

	dataLenWithFlags := uint32(len(data)) | flags
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, dataLenWithFlags)

	complete := append(append(header, data...), crcBytes...)

	out := make([]byte, len(rgb))
	copy(out, rgb)

	embedBytes2bitAtSlot(out, 0, header)
	startSlot := ((HeaderLength + IntegrityHashLen) * 8) / 2
	if scatterEnabled {
		if err := embedScatteredBytes2bit(out, startSlot, complete[HeaderLength:], password); err != nil {
			return nil, nil, err
		}
	} else {
		embedBytes2bitAtSlot(out, startSlot, complete[HeaderLength:])
	}

	integrity, err := embeddedPixelHash(out, width, height)
	if err != nil {
		return nil, nil, err
	}
	integritySlotStart := (HeaderLength * 8) / 2
	embedBytes2bitAtSlot(out, integritySlotStart, integrity)

	return out, integrity, nil
}
