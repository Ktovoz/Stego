package engine

import (
	"bufio"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func LoadImageRGB(path string) ([]byte, int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() { _ = f.Close() }()

	ext := strings.ToLower(filepath.Ext(path))
	var img image.Image
	switch ext {
	case ".png":
		img, err = png.Decode(f)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	default:
		img, _, err = image.Decode(f)
	}
	if err != nil {
		return nil, 0, 0, err
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return nil, 0, 0, errors.New("invalid image size")
	}

	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

	rgb := make([]byte, w*h*3)
	dst := 0
	for y := 0; y < h; y++ {
		row := rgba.Pix[y*rgba.Stride : y*rgba.Stride+w*4]
		for x := 0; x < w; x++ {
			rgb[dst] = row[x*4]
			rgb[dst+1] = row[x*4+1]
			rgb[dst+2] = row[x*4+2]
			dst += 3
		}
	}
	return rgb, w, h, nil
}

func SaveRGBAsPNG(path string, rgb []byte, width, height int) error {
	if len(rgb) != width*height*3 {
		return errors.New("invalid rgb buffer size")
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	pix := img.Pix
	for src, dst := 0, 0; src < len(rgb); src, dst = src+3, dst+4 {
		pix[dst] = rgb[src]
		pix[dst+1] = rgb[src+1]
		pix[dst+2] = rgb[src+2]
		pix[dst+3] = 0xFF
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	w := bufio.NewWriterSize(f, 1<<20)
	defer func() { _ = w.Flush() }()

	enc := png.Encoder{CompressionLevel: png.BestSpeed}
	return enc.Encode(w, img)
}
