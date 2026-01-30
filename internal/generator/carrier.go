package generator

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
)

type Result struct {
	PNG    []byte
	Width  int
	Height int
}

func generateTextureImage(width, height int, rng *rand.Rand) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	palette := make([][3]float32, 4)
	for i := 0; i < 4; i++ {
		palette[i] = [3]float32{
			float32(rng.Intn(256)),
			float32(rng.Intn(256)),
			float32(rng.Intn(256)),
		}
	}
	c0, c1, c2, c3 := palette[0], palette[1], palette[2], palette[3]

	invWidth := 1.0 / float32(width-1)
	invHeight := 1.0 / float32(height-1)

	angle := rng.Float64() * 2 * math.Pi
	cosA := float32(math.Cos(angle))
	sinA := float32(math.Sin(angle))

	scaleX := rng.Float32()*0.6 + 0.6
	scaleY := rng.Float32()*0.6 + 0.6
	phaseShift := rng.Float32() * 2 * math.Pi

	gridW := width / 128
	if gridW < 24 {
		gridW = 24
	}
	gridH := height / 128
	if gridH < 24 {
		gridH = 24
	}

	lowFreqNoise := make([][]float32, gridH)
	for y := 0; y < gridH; y++ {
		lowFreqNoise[y] = make([]float32, gridW)
		for x := 0; x < gridW; x++ {
			u1 := rng.Float64()
			u2 := rng.Float64()
			z0 := float32(math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2))
			lowFreqNoise[y][x] = z0
		}
	}

	minNoise, maxNoise := lowFreqNoise[0][0], lowFreqNoise[0][0]
	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			if lowFreqNoise[y][x] < minNoise {
				minNoise = lowFreqNoise[y][x]
			}
			if lowFreqNoise[y][x] > maxNoise {
				maxNoise = lowFreqNoise[y][x]
			}
		}
	}
	rangeNoise := maxNoise - minNoise
	if rangeNoise < 1e-6 {
		rangeNoise = 1e-6
	}
	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			lowFreqNoise[y][x] = (lowFreqNoise[y][x] - minNoise) / rangeNoise
		}
	}

	cx := rng.Float32()*0.3 + 0.35
	cy := rng.Float32()*0.3 + 0.35
	vignetteStrength := rng.Float32()*0.07 + 0.05

	cloudStrength := rng.Float32()*8.0 + 10.0

	gridScaleX := float64(gridW) / float64(width)
	gridScaleY := float64(gridH) / float64(height)

	for y := 0; y < height; y++ {
		yy := float32(y) * invHeight
		dy := yy - cy
		dy2 := dy * dy

		for x := 0; x < width; x++ {
			xx := float32(x) * invWidth

			t := cosA*xx + sinA*yy
			if t < 0 {
				t = 0
			} else if t > 1 {
				t = 1
			}

			invT := 1 - t
			base := [3]float32{
				invT*c0[0] + t*c1[0],
				invT*c0[1] + t*c1[1],
				invT*c0[2] + t*c1[2],
			}

			arg := 2.0*math.Pi*(float64(xx*scaleX)+float64(yy*scaleY)) + float64(phaseShift)
			t2 := 0.5 + 0.5*float32(math.Sin(arg))
			if t2 < 0 {
				t2 = 0
			} else if t2 > 1 {
				t2 = 1
			}

			invT2 := 1 - t2
			layer2 := [3]float32{
				invT2*c2[0] + t2*c3[0],
				invT2*c2[1] + t2*c3[1],
				invT2*c2[2] + t2*c3[2],
			}

			base = [3]float32{
				base[0]*0.75 + layer2[0]*0.25,
				base[1]*0.75 + layer2[1]*0.25,
				base[2]*0.75 + layer2[2]*0.25,
			}

			srcX := float64(x) * gridScaleX
			srcY := float64(y) * gridScaleY
			x0 := int(srcX)
			y0 := int(srcY)
			x1 := x0 + 1
			y1 := y0 + 1
			if x1 >= gridW {
				x1 = gridW - 1
			}
			if y1 >= gridH {
				y1 = gridH - 1
			}
			fx := float32(srcX - float64(x0))
			fy := float32(srcY - float64(y0))

			n00 := lowFreqNoise[y0][x0]
			n01 := lowFreqNoise[y0][x1]
			n10 := lowFreqNoise[y1][x0]
			n11 := lowFreqNoise[y1][x1]
			noise := (1-fx)*(1-fy)*n00 + fx*(1-fy)*n01 + (1-fx)*fy*n10 + fx*fy*n11

			cloud := (noise - 0.5) * cloudStrength
			base = [3]float32{
				base[0] + cloud,
				base[1] + cloud,
				base[2] + cloud,
			}

			dx := xx - cx
			r2 := dx*dx + dy2
			vignette := 1.0 - r2/0.9
			if vignette < 0 {
				vignette = 0
			}
			vignette *= vignetteStrength
			vignette = 1.0 - vignette
			if vignette < 0 {
				vignette = 0
			}

			base = [3]float32{
				base[0] * vignette,
				base[1] * vignette,
				base[2] * vignette,
			}

			clamp := func(v float32) uint8 {
				if v < 0 {
					return 0
				}
				if v > 255 {
					return 255
				}
				return uint8(v)
			}

			img.SetRGBA(x, y, color.RGBA{
				R: clamp(base[0]),
				G: clamp(base[1]),
				B: clamp(base[2]),
				A: 0xFF,
			})
		}
	}

	blurRadius := 0.8
	if width > 2000 || height > 2000 {
		blurRadius = 0.3
	}
	if width > 3000 || height > 3000 {
		blurRadius = 0
	}
	if blurRadius > 0 {
		img = gaussianBlur(img, blurRadius)
	}

	return img
}

func gaussianBlur(img *image.RGBA, radius float64) *image.RGBA {
	if radius < 0.5 {
		return img
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)

	kernelSize := int(radius*3) + 1
	if kernelSize < 3 {
		kernelSize = 3
	}
	if kernelSize%2 == 0 {
		kernelSize++
	}

	kernel := make([]float64, kernelSize)
	sum := 0.0
	sigma := radius / 3.0
	center := kernelSize / 2
	for i := 0; i < kernelSize; i++ {
		x := float64(i - center)
		kernel[i] = math.Exp(-(x * x) / (2 * sigma * sigma))
		sum += kernel[i]
	}
	for i := range kernel {
		kernel[i] /= sum
	}

	temp := make([][]float64, h)
	for y := 0; y < h; y++ {
		temp[y] = make([]float64, w*4)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sumR, sumG, sumB, sumA float64
			for k := 0; k < kernelSize; k++ {
				xk := x - center + k
				if xk < 0 {
					xk = 0
				} else if xk >= w {
					xk = w - 1
				}
				r, g, b, a := img.At(xk, y).RGBA()
				weight := kernel[k]
				sumR += float64(r>>8) * weight
				sumG += float64(g>>8) * weight
				sumB += float64(b>>8) * weight
				sumA += float64(a>>8) * weight
			}
			temp[y][x*4+0] = sumR
			temp[y][x*4+1] = sumG
			temp[y][x*4+2] = sumB
			temp[y][x*4+3] = sumA
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sumR, sumG, sumB, sumA float64
			for k := 0; k < kernelSize; k++ {
				yk := y - center + k
				if yk < 0 {
					yk = 0
				} else if yk >= h {
					yk = h - 1
				}
				weight := kernel[k]
				sumR += temp[yk][x*4+0] * weight
				sumG += temp[yk][x*4+1] * weight
				sumB += temp[yk][x*4+2] * weight
				sumA += temp[yk][x*4+3] * weight
			}

			result.SetRGBA(x, y, color.RGBA{
				R: uint8(sumR + 0.5),
				G: uint8(sumG + 0.5),
				B: uint8(sumB + 0.5),
				A: uint8(sumA + 0.5),
			})
		}
	}

	return result
}

func GenerateCarrierPNG(targetBytes int64, seed int64, noise bool) (Result, error) {
	if targetBytes <= 0 {
		return Result{}, errors.New("target bytes must be > 0")
	}
	if seed == 0 {
		seed = 1
	}
	rng := rand.New(rand.NewSource(seed))

	requiredPixels := math.Ceil((float64(targetBytes) + 32) / 0.75)
	side := int(math.Ceil(math.Sqrt(requiredPixels)))

	if side < 64 {
		side = 64
	}

	img := generateTextureImage(side, side, rng)

	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestSpeed}
	if err := enc.Encode(&buf, img); err != nil {
		return Result{}, err
	}

	return Result{PNG: buf.Bytes(), Width: side, Height: side}, nil
}
