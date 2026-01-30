package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	eccMagic = []byte("RS1")
)

const (
	RSK    = 223
	RSNSym = 32
)

func ECCWrapRS(data []byte) ([]byte, error) {
	framed := make([]byte, 4+len(data))
	binary.LittleEndian.PutUint32(framed[0:4], uint32(len(data)))
	copy(framed[4:], data)

	blocks := (len(framed) + RSK - 1) / RSK
	cwLen := RSK + RSNSym

	codewords := make([][]byte, blocks)
	for i := 0; i < blocks; i++ {
		chunk := framed[i*RSK:]
		if len(chunk) > RSK {
			chunk = chunk[:RSK]
		}
		if len(chunk) < RSK {
			padded := make([]byte, RSK)
			copy(padded, chunk)
			chunk = padded
		}
		cw := rsEncode(chunk, RSNSym)
		codewords[i] = cw
	}

	payload := rsInterleave(codewords, cwLen)
	header := make([]byte, 3+2+2+4)
	copy(header[0:3], eccMagic)
	binary.LittleEndian.PutUint16(header[3:5], uint16(RSK))
	binary.LittleEndian.PutUint16(header[5:7], uint16(RSNSym))
	binary.LittleEndian.PutUint32(header[7:11], uint32(len(framed)))

	return append(header, payload...), nil
}

func ECCUnwrapRS(blob []byte) ([]byte, error) {
	if !bytes.HasPrefix(blob, eccMagic) {
		return blob, nil
	}
	if len(blob) < 11 {
		return nil, errors.New("ecc header corrupted")
	}
	k := int(binary.LittleEndian.Uint16(blob[3:5]))
	nsym := int(binary.LittleEndian.Uint16(blob[5:7]))
	framedLen := int(binary.LittleEndian.Uint32(blob[7:11]))
	if k != RSK || nsym != RSNSym {
		return nil, errors.New("unsupported ecc parameters")
	}
	cwLen := k + nsym
	interleaved := blob[11:]
	if len(interleaved)%cwLen != 0 {
		return nil, errors.New("ecc payload length invalid")
	}
	blocks := len(interleaved) / cwLen
	codewords := rsDeinterleave(interleaved, blocks, cwLen)

	decoded := make([]byte, 0, blocks*k)
	for _, cw := range codewords {
		msg, err := rsDecode(cw, k, nsym)
		if err != nil {
			return nil, err
		}
		decoded = append(decoded, msg...)
	}
	if framedLen < 4 || framedLen > len(decoded) {
		return nil, errors.New("ecc decoded length invalid")
	}
	framed := decoded[:framedLen]
	n := int(binary.LittleEndian.Uint32(framed[0:4]))
	if n != framedLen-4 {
		return nil, errors.New("ecc length mismatch")
	}
	out := make([]byte, n)
	copy(out, framed[4:])
	return out, nil
}

func rsInterleave(codewords [][]byte, cwLen int) []byte {
	blocks := len(codewords)
	out := make([]byte, 0, blocks*cwLen)
	for col := 0; col < cwLen; col++ {
		for row := 0; row < blocks; row++ {
			out = append(out, codewords[row][col])
		}
	}
	return out
}

func rsDeinterleave(interleaved []byte, blocks, cwLen int) [][]byte {
	out := make([][]byte, blocks)
	for i := range out {
		out[i] = make([]byte, cwLen)
	}
	idx := 0
	for col := 0; col < cwLen; col++ {
		for row := 0; row < blocks; row++ {
			out[row][col] = interleaved[idx]
			idx++
		}
	}
	return out
}
