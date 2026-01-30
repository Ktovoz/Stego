package crypto

import "errors"

var (
	gfExp [512]byte
	gfLog [256]byte
)

func init() {
	x := byte(1)
	for i := 0; i < 255; i++ {
		gfExp[i] = x
		gfLog[x] = byte(i)
		x = gfMulNoLUT(x, 2)
	}
	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}
}

func gfMulNoLUT(x, y byte) byte {
	var r byte
	for y > 0 {
		if (y & 1) != 0 {
			r ^= x
		}
		y >>= 1
		carry := (x & 0x80) != 0
		x <<= 1
		if carry {
			x ^= 0x1d
		}
	}
	return r
}

func gfMul(x, y byte) byte {
	if x == 0 || y == 0 {
		return 0
	}
	return gfExp[int(gfLog[x])+int(gfLog[y])]
}

func gfDiv(x, y byte) byte {
	if y == 0 {
		panic("gf divide by zero")
	}
	if x == 0 {
		return 0
	}
	return gfExp[int(gfLog[x])+255-int(gfLog[y])]
}

func gfPow2(power int) byte {
	if power == 0 {
		return 1
	}
	logx := int(gfLog[2])
	return gfExp[(logx*power)%255]
}

func polyScale(p []byte, x byte) []byte {
	out := make([]byte, len(p))
	for i := range p {
		out[i] = gfMul(p[i], x)
	}
	return out
}

func polyAdd(p, q []byte) []byte {
	if len(p) < len(q) {
		p, q = q, p
	}
	out := make([]byte, len(p))
	copy(out, p)
	offset := len(p) - len(q)
	for i := range q {
		out[i+offset] ^= q[i]
	}
	return out
}

func polyMul(p, q []byte) []byte {
	out := make([]byte, len(p)+len(q)-1)
	for j := 0; j < len(q); j++ {
		for i := 0; i < len(p); i++ {
			out[i+j] ^= gfMul(p[i], q[j])
		}
	}
	return out
}

func polyEval(poly []byte, x byte) byte {
	var y byte
	for i := 0; i < len(poly); i++ {
		y = gfMul(y, x) ^ poly[i]
	}
	return y
}

func rsGeneratorPoly(nsym int) []byte {
	g := []byte{1}
	for i := 0; i < nsym; i++ {
		g = polyMul(g, []byte{1, gfPow2(i + 1)})
	}
	return g
}

func rsEncode(msg []byte, nsym int) []byte {
	gen := rsGeneratorPoly(nsym)
	out := make([]byte, len(msg)+nsym)
	copy(out, msg)

	for i := 0; i < len(msg); i++ {
		coef := out[i]
		if coef == 0 {
			continue
		}
		for j := 1; j < len(gen); j++ {
			out[i+j] ^= gfMul(gen[j], coef)
		}
	}
	copy(out, msg)
	return out
}

func rsCalcSyndromes(msg []byte, nsym int) []byte {
	synd := make([]byte, nsym+1)
	synd[0] = 0
	for i := 0; i < nsym; i++ {
		synd[i+1] = polyEval(msg, gfPow2(i+1))
	}
	return synd
}

func rsCheck(synd []byte) bool {
	for i := 1; i < len(synd); i++ {
		if synd[i] != 0 {
			return false
		}
	}
	return true
}

func rsFindErrorLocator(synd []byte, nsym int) ([]byte, error) {
	errLoc := []byte{1}
	oldLoc := []byte{1}

	for i := 0; i < nsym; i++ {
		delta := synd[i+1]
		for j := 1; j < len(errLoc); j++ {
			delta ^= gfMul(errLoc[len(errLoc)-1-j], synd[i+1-j])
		}

		oldLoc = append(oldLoc, 0)
		if delta != 0 {
			if len(oldLoc) > len(errLoc) {
				newLoc := polyScale(oldLoc, delta)
				oldLoc = polyScale(errLoc, gfDiv(1, delta))
				errLoc = newLoc
			}
			errLoc = polyAdd(errLoc, polyScale(oldLoc, delta))
		}
	}

	for len(errLoc) > 0 && errLoc[0] == 0 {
		errLoc = errLoc[1:]
	}
	errCount := len(errLoc) - 1
	if errCount*2 > nsym {
		return nil, errors.New("too many errors to correct")
	}
	return errLoc, nil
}

func rsFindErrors(errLoc []byte, nmess int) ([]int, error) {
	errs := len(errLoc) - 1
	if errs == 0 {
		return nil, nil
	}
	errPos := make([]int, 0, errs)
	loc := make([]byte, len(errLoc))
	for i := range errLoc {
		loc[i] = errLoc[len(errLoc)-1-i]
	}
	for i := 0; i < nmess; i++ {
		if polyEval(loc, gfPow2(i)) == 0 {
			errPos = append(errPos, nmess-1-i)
		}
	}
	if len(errPos) != errs {
		return nil, errors.New("could not locate errors")
	}
	return errPos, nil
}

func rsErrorEvaluator(synd, errLoc []byte, nsym int) []byte {
	product := polyMul(synd, errLoc)
	if len(product) <= nsym {
		return product
	}
	return product[len(product)-nsym:]
}

func polyDeriv(p []byte) []byte {
	deg := len(p) - 1
	if deg <= 0 {
		return []byte{0}
	}
	out := make([]byte, 0, len(p)-1)
	for i := 0; i < len(p)-1; i++ {
		power := deg - i
		if power%2 == 1 {
			out = append(out, p[i])
		}
	}
	if len(out) == 0 {
		return []byte{0}
	}
	return out
}

func rsDecode(codeword []byte, k, nsym int) ([]byte, error) {
	if len(codeword) != k+nsym {
		return nil, errors.New("invalid codeword length")
	}
	synd := rsCalcSyndromes(codeword, nsym)
	if rsCheck(synd) {
		msg := make([]byte, k)
		copy(msg, codeword[:k])
		return msg, nil
	}
	errLoc, err := rsFindErrorLocator(synd, nsym)
	if err != nil {
		return nil, err
	}
	errPos, err := rsFindErrors(errLoc, len(codeword))
	if err != nil {
		return nil, err
	}
	corrected, err := rsCorrect(codeword, synd, errLoc, errPos, nsym)
	if err != nil {
		return nil, err
	}
	synd2 := rsCalcSyndromes(corrected, nsym)
	if !rsCheck(synd2) {
		return nil, errors.New("could not correct message")
	}
	msg := make([]byte, k)
	copy(msg, corrected[:k])
	return msg, nil
}

func rsCorrect(msg []byte, synd []byte, errLoc []byte, errPos []int, nsym int) ([]byte, error) {
	nmess := len(msg)
	errEval := rsErrorEvaluator(synd, errLoc, nsym)
	errLocDeriv := polyDeriv(errLoc)
	if len(errLocDeriv) == 0 || (len(errLocDeriv) == 1 && errLocDeriv[0] == 0) {
		return nil, errors.New("invalid error locator derivative")
	}

	out := make([]byte, len(msg))
	copy(out, msg)
	for _, p := range errPos {
		coefPos := nmess - 1 - p
		x := gfPow2(coefPos + 1)
		xInv := gfDiv(1, x)

		y := polyEval(errEval, xInv)
		d := polyEval(errLocDeriv, xInv)
		if d == 0 {
			return nil, errors.New("division by zero during correction")
		}
		mag := gfDiv(y, d)
		out[p] ^= mag
	}
	return out, nil
}
