package app

import (
	"context"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"

	"stego/internal/engine"
)

func selectCarrierImage(ctx context.Context, eng *engine.Engine, carrierDir string, requiredBytes int, preferLargest bool) (string, error) {
	if err := os.MkdirAll(carrierDir, 0o755); err != nil {
		return "", err
	}

	entries, err := os.ReadDir(carrierDir)
	if err != nil {
		return "", err
	}
	type cand struct {
		path     string
		capacity int
		score    float64
	}
	var best *cand
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".bmp" && ext != ".webp" && ext != ".tiff" && ext != ".gif" {
			continue
		}
		path := filepath.Join(carrierDir, e.Name())
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		cfg, _, err := image.DecodeConfig(f)
		_ = f.Close()
		if err != nil {
			continue
		}
		capacity := eng.CalculateMaxCapacity(cfg.Width, cfg.Height, false)
		if requiredBytes > capacity {
			continue
		}
		score := 0.0
		if !preferLargest {
			score = quickTextureScore(path, 256)
		}
		c := cand{path: path, capacity: capacity, score: score}
		if best == nil {
			best = &c
			continue
		}
		if preferLargest {
			if c.capacity > best.capacity {
				best = &c
			}
			continue
		}
		if c.score > best.score {
			best = &c
		}
	}
	if best == nil {
		return "", errors.New("no suitable carrier image found")
	}
	return best.path, nil
}

func quickTextureScore(imagePath string, sampleSize int) float64 {
	f, err := os.Open(imagePath)
	if err != nil {
		return -1
	}
	img, _, err := image.Decode(f)
	_ = f.Close()
	if err != nil {
		return -1
	}
	if sampleSize <= 0 {
		sampleSize = 256
	}
	gray := image.NewGray(image.Rect(0, 0, sampleSize, sampleSize))
	draw.ApproxBiLinear.Scale(gray, gray.Bounds(), img, img.Bounds(), draw.Over, nil)

	var sum float64
	w, h := gray.Bounds().Dx(), gray.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := gray.GrayAt(x, y).Y
			l := gray.GrayAt(maxInt(0, x-1), y).Y
			r := gray.GrayAt(minInt(w-1, x+1), y).Y
			u := gray.GrayAt(x, maxInt(0, y-1)).Y
			d := gray.GrayAt(x, minInt(h-1, y+1)).Y
			lap := int(l) + int(r) + int(u) + int(d) - 4*int(c)
			sum += math.Abs(float64(lap))
		}
	}
	return sum / float64(w*h)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
