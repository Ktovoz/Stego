package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"stego/internal/generator"
	"stego/internal/models"
)

func RunGenerateCarrier(ctx context.Context, req models.GenerateRequest, emit func(models.ProgressEvent), taskID string, logf PerfLogger) error {
	if emit == nil {
		emit = func(models.ProgressEvent) {}
	}
	startAll := time.Now()
	ok := false
	defer func() {
		logPerf(logf, "generate", taskID, "Total", time.Since(startAll), fmt.Sprintf("ok=%t count=%d targetBytes=%d", ok, req.Count, req.TargetBytes))
	}()

	outDir := strings.TrimSpace(req.OutputDir)
	if outDir == "" {
		outDir = "./images"
	}
	if req.Count <= 0 {
		return errors.New("count must be > 0")
	}
	if req.TargetBytes <= 0 {
		return errors.New("targetBytes must be > 0")
	}
	prefix := strings.TrimSpace(req.Prefix)
	if prefix == "" {
		prefix = "carrier"
	}
	seedBase := req.RandomSeed
	if seedBase == 0 {
		seedBase = time.Now().UnixNano()
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	var genTotal time.Duration
	var writeTotal time.Duration
	for i := 1; i <= req.Count; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		progress := int(float64(i-1) / float64(req.Count) * 100)
		emit(models.ProgressEvent{Progress: progress, Current: i - 1, Total: req.Count, Message: "生成图片..."})

		t0 := time.Now()
		res, err := generator.GenerateCarrierPNG(req.TargetBytes, seedBase+int64(i), req.NoiseEnabled)
		if err != nil {
			emit(models.ProgressEvent{Progress: progress, Error: err.Error(), Done: true})
			return err
		}
		genTotal += time.Since(t0)
		name := prefix + "_" + strconv.Itoa(i) + ".png"
		path := filepath.Join(outDir, name)
		t0 = time.Now()
		if err := os.WriteFile(path, res.PNG, 0o644); err != nil {
			return err
		}
		writeTotal += time.Since(t0)
	}

	emit(models.ProgressEvent{Progress: 100, Current: req.Count, Total: req.Count, Message: "完成", Done: true})
	if req.Count > 0 {
		logPerf(logf, "generate", taskID, "GenerateTotal", genTotal, fmt.Sprintf("count=%d avg=%s", req.Count, formatDuration(genTotal/time.Duration(req.Count))))
		logPerf(logf, "generate", taskID, "WriteTotal", writeTotal, fmt.Sprintf("count=%d avg=%s", req.Count, formatDuration(writeTotal/time.Duration(req.Count))))
	}
	ok = true
	return nil
}
