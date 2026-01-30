package app

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"stego/internal/config"
	"stego/internal/crypto"
	"stego/internal/engine"
	"stego/internal/models"
)

func RunDecrypt(ctx context.Context, cfg map[string]string, req models.DecryptRequest, emit func(models.ProgressEvent), taskID string, logf PerfLogger) error {
	if emit == nil {
		emit = func(models.ProgressEvent) {}
	}
	startAll := time.Now()
	ok := false
	defer func() {
		logPerf(logf, "decrypt", taskID, "Total", time.Since(startAll), fmt.Sprintf("ok=%t", ok))
	}()
	password := strings.TrimSpace(req.Password)
	if password == "" {
		password = cfg[config.KeyDefaultDecryptPassword]
	}
	outputDir := strings.TrimSpace(req.OutputDir)
	if outputDir == "" {
		outputDir = cfg[config.KeyDefaultOutputDir]
	}
	identifier := strings.TrimSpace(req.Identifier)
	if identifier == "" {
		identifier = "stego"
	}

	emit(models.ProgressEvent{Progress: 0, Message: "读取图片..."})
	t0 := time.Now()
	rgb, w, h, err := engine.LoadImageRGB(req.ImagePath)
	if err != nil {
		emit(models.ProgressEvent{Progress: 0, Error: err.Error(), Done: true})
		return err
	}
	logPerf(logf, "decrypt", taskID, "LoadImage", time.Since(t0), fmt.Sprintf("w=%d h=%d", w, h))
	if err := ctx.Err(); err != nil {
		return err
	}

	emit(models.ProgressEvent{Progress: 20, Message: "提取数据..."})
	eng := engine.New(1024 * 1024)
	t0 = time.Now()
	extracted, _, _, _, err := eng.Extract(rgb, w, h, password)
	if err != nil {
		emit(models.ProgressEvent{Progress: 20, Error: err.Error(), Done: true})
		return err
	}
	logPerf(logf, "decrypt", taskID, "Extract", time.Since(t0), fmt.Sprintf("bytes=%d", len(extracted)))

	emit(models.ProgressEvent{Progress: 40, Message: "纠错解码..."})
	t0 = time.Now()
	extracted, err = crypto.ECCUnwrapRS(extracted)
	if err != nil {
		emit(models.ProgressEvent{Progress: 40, Error: err.Error(), Done: true})
		return err
	}
	logPerf(logf, "decrypt", taskID, "ECCUnwrap", time.Since(t0), fmt.Sprintf("bytes=%d", len(extracted)))
	if len(extracted) < engine.MetadataLengthSize {
		return errors.New("data format invalid: metadata length missing")
	}
	metaLen := int(binary.LittleEndian.Uint32(extracted[:engine.MetadataLengthSize]))
	metaEnd := engine.MetadataLengthSize + metaLen
	if metaLen < 0 || metaEnd > len(extracted) {
		return errors.New("data format invalid: metadata length out of range")
	}

	var meta encryptMetadata
	if err := json.Unmarshal(extracted[engine.MetadataLengthSize:metaEnd], &meta); err != nil {
		return err
	}
	encrypted := extracted[metaEnd:]
	minSize := meta.SaltLength + meta.NonceLength + meta.TagLength
	if len(encrypted) < minSize {
		return errors.New("encrypted payload incomplete")
	}
	salt := encrypted[:meta.SaltLength]
	nonce := encrypted[meta.SaltLength : meta.SaltLength+meta.NonceLength]
	tag := encrypted[meta.SaltLength+meta.NonceLength : meta.SaltLength+meta.NonceLength+meta.TagLength]
	ciphertext := encrypted[meta.SaltLength+meta.NonceLength+meta.TagLength:]

	emit(models.ProgressEvent{Progress: 60, Message: "解密..."})
	t0 = time.Now()
	key := crypto.PBKDF2Compat(password, salt, meta.PBKDF2Iterations, meta.KeyLength)
	logPerf(logf, "decrypt", taskID, "KDF", time.Since(t0), fmt.Sprintf("iters=%d keyLen=%d", meta.PBKDF2Iterations, meta.KeyLength))
	t0 = time.Now()
	plain, err := crypto.DecryptAESGCM(key, nonce, ciphertext, tag)
	if err != nil {
		emit(models.ProgressEvent{Progress: 60, Error: err.Error(), Done: true})
		return err
	}
	logPerf(logf, "decrypt", taskID, "Decrypt", time.Since(t0), fmt.Sprintf("plainBytes=%d", len(plain)))

	outBase := filepath.Join(outputDir, "extracted")
	if err := os.MkdirAll(outBase, 0o755); err != nil {
		return err
	}
	emit(models.ProgressEvent{Progress: 80, Message: "写出文件..."})
	t0 = time.Now()
	if isZip(plain) {
		dest := filepath.Join(outBase, identifier+"_"+filepath.Base(strings.TrimSuffix(req.ImagePath, filepath.Ext(req.ImagePath))))
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return err
		}
		if err := unzipToDir(plain, dest); err != nil {
			return err
		}
	} else {
		outFile := filepath.Join(outBase, filepath.Base(req.ImagePath)+"_extracted.bin")
		if err := os.WriteFile(outFile, plain, 0o644); err != nil {
			return err
		}
	}
	logPerf(logf, "decrypt", taskID, "WriteOutput", time.Since(t0), "")

	emit(models.ProgressEvent{Progress: 100, Message: "完成", Done: true})
	ok = true
	return nil
}
