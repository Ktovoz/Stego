package engine

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const (
	HeaderLength       = 4
	CRCLength          = 4
	IntegrityHashLen   = 16
	MetadataLengthSize = 4

	IntegrityFlag = 0x80000000
	ScatterFlag   = 0x40000000
)

type Engine struct {
	ChunkSize int
}

func New(chunkSize int) *Engine {
	if chunkSize <= 0 {
		chunkSize = 1024 * 1024
	}
	return &Engine{ChunkSize: chunkSize}
}

func (e *Engine) CalculateMaxCapacity(width, height int, includeOverhead bool) int {
	if width <= 0 || height <= 0 {
		return 0
	}
	base := (width * height * 3 * 2) / 8
	if includeOverhead {
		return base - 32
	}
	return base
}

func calculateCRC32(data []byte) []byte {
	out := make([]byte, 4)
	binary.LittleEndian.PutUint32(out, crc32.ChecksumIEEE(data))
	return out
}

func verifyCRC32(data []byte, crcBytes []byte) bool {
	if len(crcBytes) != 4 {
		return false
	}
	expected := binary.LittleEndian.Uint32(crcBytes)
	actual := crc32.ChecksumIEEE(data)
	return expected == actual
}

func embeddedPixelHash(rgb []byte, width, height int) ([]byte, error) {
	if len(rgb) != width*height*3 {
		return nil, errors.New("invalid rgb buffer size")
	}

	hashBitOffset := HeaderLength * 8
	hashSlotStart := hashBitOffset / 2
	hashSlots := (IntegrityHashLen * 8) / 2
	hashSlotEnd := hashSlotStart + hashSlots
	if hashSlotEnd > len(rgb) {
		hashSlotEnd = len(rgb)
	}

	header := make([]byte, 8)
	binary.LittleEndian.PutUint32(header[0:4], uint32(width))
	binary.LittleEndian.PutUint32(header[4:8], uint32(height))

	h := sha256.New()
	_, _ = h.Write(header)
	if hashSlotStart > 0 {
		_, _ = h.Write(rgb[:hashSlotStart])
	}
	if hashSlotEnd > hashSlotStart {
		buf := make([]byte, 32*1024)
		for i := hashSlotStart; i < hashSlotEnd; {
			n := hashSlotEnd - i
			if n > len(buf) {
				n = len(buf)
			}
			for j := 0; j < n; j++ {
				buf[j] = rgb[i+j] & 0xFC
			}
			_, _ = h.Write(buf[:n])
			i += n
		}
	}
	if hashSlotEnd < len(rgb) {
		_, _ = h.Write(rgb[hashSlotEnd:])
	}
	sum := h.Sum(nil)
	out := make([]byte, IntegrityHashLen)
	copy(out, sum[:IntegrityHashLen])
	return out, nil
}
