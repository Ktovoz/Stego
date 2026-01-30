package app

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"stego/internal/config"
	"stego/internal/crypto"
	"stego/internal/engine"
	"stego/internal/models"
)

type encryptMetadata struct {
	Algorithm        string `json:"algorithm"`
	KeyLength        int    `json:"key_length"`
	SaltLength       int    `json:"salt_length"`
	NonceLength      int    `json:"nonce_length"`
	TagLength        int    `json:"tag_length"`
	PBKDF2Iterations int    `json:"pbkdf2_iterations"`
}

func RunEncrypt(ctx context.Context, cfg map[string]string, req models.EncryptRequest, emit func(models.ProgressEvent), taskID string, logf PerfLogger) error {
	if emit == nil {
		emit = func(models.ProgressEvent) {}
	}
	startAll := time.Now()
	ok := false
	defer func() {
		logPerf(logf, "encrypt", taskID, "Total", time.Since(startAll), fmt.Sprintf("ok=%t", ok))
	}()
	password := strings.TrimSpace(req.Password)
	if password == "" {
		password = cfg[config.KeyDefaultEncryptPassword]
	}
	carrierDir := strings.TrimSpace(req.CarrierDir)
	if carrierDir == "" {
		carrierDir = cfg[config.KeyDefaultCarrierDir]
	}
	outputDir := strings.TrimSpace(req.OutputDir)
	if outputDir == "" {
		outputDir = cfg[config.KeyDefaultOutputDir]
	}
	outputFileName := strings.TrimSpace(req.OutputFileName)
	if outputFileName == "" {
		outputFileName = cfg[config.KeyDefaultEncryptOutputName]
		if outputFileName == "" {
			outputFileName = "encrypted"
		}
	}
	scatter := true
	if req.Scatter != nil {
		scatter = *req.Scatter
	}

	emit(models.ProgressEvent{Progress: 0, Message: "读取数据源..."})
	t0 := time.Now()
	data, _, err := readDataSource(ctx, req.DataSourcePath)
	if err != nil {
		emit(models.ProgressEvent{Progress: 0, Error: err.Error(), Done: true})
		return err
	}
	logPerf(logf, "encrypt", taskID, "ReadDataSource", time.Since(t0), fmt.Sprintf("bytes=%d", len(data)))
	if err := ctx.Err(); err != nil {
		return err
	}

	cryptoCfg := crypto.DefaultAESGCMConfig()
	metaProbe := encryptMetadata{
		Algorithm:        "AES-GCM",
		KeyLength:        cryptoCfg.KeyLength,
		SaltLength:       cryptoCfg.SaltLength,
		NonceLength:      cryptoCfg.NonceLen,
		TagLength:        cryptoCfg.TagLen,
		PBKDF2Iterations: cryptoCfg.Iterations,
	}
	metaJSON, _ := json.Marshal(metaProbe)
	requiredPayloadBytes := estimateRequiredPayloadBytes(int64(len(data)), int64(len(metaJSON)), cryptoCfg.SaltLength, cryptoCfg.NonceLen, cryptoCfg.TagLen)
	requiredBytesInCarrier := engine.HeaderLength + engine.IntegrityHashLen + int(requiredPayloadBytes) + engine.CRCLength

	emit(models.ProgressEvent{Progress: 10, Message: "选择载体图片..."})
	carrierPath := strings.TrimSpace(req.CarrierImagePath)
	eng := engine.New(1024 * 1024)
	if carrierPath == "" {
		t0 = time.Now()
		p, err := selectCarrierImage(ctx, eng, carrierDir, requiredBytesInCarrier, req.PreferLargestImage)
		if err != nil {
			emit(models.ProgressEvent{Progress: 10, Error: err.Error(), Done: true})
			return err
		}
		carrierPath = p
		logPerf(logf, "encrypt", taskID, "SelectCarrierImage", time.Since(t0), "")
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	emit(models.ProgressEvent{Progress: 20, Message: "加密并纠错编码..."})

	t0 = time.Now()
	salt, err := crypto.RandomBytes(cryptoCfg.SaltLength)
	if err != nil {
		return err
	}
	nonce, err := crypto.RandomBytes(cryptoCfg.NonceLen)
	if err != nil {
		return err
	}
	key := crypto.PBKDF2Compat(password, salt, cryptoCfg.Iterations, cryptoCfg.KeyLength)
	ciphertext, tag, err := crypto.EncryptAESGCM(key, nonce, data)
	if err != nil {
		return err
	}

	meta := encryptMetadata{
		Algorithm:        "AES-GCM",
		KeyLength:        cryptoCfg.KeyLength,
		SaltLength:       cryptoCfg.SaltLength,
		NonceLength:      cryptoCfg.NonceLen,
		TagLength:        cryptoCfg.TagLen,
		PBKDF2Iterations: cryptoCfg.Iterations,
	}
	metaJSON, err = json.Marshal(meta)
	if err != nil {
		return err
	}
	metaLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(metaLen, uint32(len(metaJSON)))
	fullData := append(append(append(append(metaLen, metaJSON...), salt...), nonce...), tag...)
	fullData = append(fullData, ciphertext...)

	wrapped, err := crypto.ECCWrapRS(fullData)
	if err != nil {
		return err
	}
	logPerf(logf, "encrypt", taskID, "Encrypt+ECCWrap", time.Since(t0), fmt.Sprintf("wrappedBytes=%d", len(wrapped)))

	emit(models.ProgressEvent{Progress: 50, Message: "嵌入数据..."})
	t0 = time.Now()
	rgb, w, h, err := engine.LoadImageRGB(carrierPath)
	if err != nil {
		return err
	}
	logPerf(logf, "encrypt", taskID, "LoadCarrierImage", time.Since(t0), fmt.Sprintf("w=%d h=%d", w, h))
	t0 = time.Now()
	outRGB, _, err := eng.Hide(rgb, w, h, wrapped, password, scatter)
	if err != nil {
		return err
	}
	logPerf(logf, "encrypt", taskID, "Hide", time.Since(t0), "")

	outFile := filepath.Join(outputDir, "encrypted", outputFileName)
	if filepath.Ext(outFile) == "" {
		outFile += ".png"
	}
	outFile = uniqueFilePath(outFile)

	if err := os.MkdirAll(filepath.Dir(outFile), 0o755); err != nil {
		return err
	}

	emit(models.ProgressEvent{Progress: 90, Message: "保存图片..."})
	t0 = time.Now()
	if err := engine.SaveRGBAsPNG(outFile, outRGB, w, h); err != nil {
		return err
	}
	logPerf(logf, "encrypt", taskID, "SavePNG", time.Since(t0), filepath.Base(outFile))

	emit(models.ProgressEvent{Progress: 100, Message: "完成", Done: true})
	ok = true
	return nil
}

func uniqueFilePath(path string) string {
	if _, err := os.Stat(path); err != nil {
		return path
	}
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for i := 2; i < 10000; i++ {
		p := base + "_" + strconv.Itoa(i) + ext
		if _, err := os.Stat(p); err != nil {
			return p
		}
	}
	return path
}

func estimateRequiredPayloadBytes(plainLen int64, metaJSONLen int64, saltLen int, nonceLen int, tagLen int) int64 {
	fullDataLen := int64(4) + metaJSONLen + int64(saltLen) + int64(nonceLen) + int64(tagLen) + plainLen
	framedLen := int64(4) + fullDataLen
	blocks := (framedLen + crypto.RSK - 1) / crypto.RSK
	wrappedLen := int64(3+2+2+4) + blocks*int64(crypto.RSK+crypto.RSNSym)
	return wrappedLen + 256
}
