package app

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var zipMagic = []byte{'P', 'K', 0x03, 0x04}

func readDataSource(ctx context.Context, dataSourcePath string) ([]byte, string, error) {
	info, err := os.Stat(dataSourcePath)
	if err != nil {
		return nil, "", err
	}
	baseName := strings.TrimSuffix(filepath.Base(dataSourcePath), filepath.Ext(dataSourcePath))
	if info.IsDir() {
		b, err := zipDirectory(ctx, dataSourcePath)
		return b, filepath.Base(dataSourcePath), err
	}
	b, err := os.ReadFile(dataSourcePath)
	return b, baseName, err
}

func zipDirectory(ctx context.Context, dir string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	defer func() { _ = zw.Close() }()

	root := filepath.Clean(dir)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		info, err := d.Info()
		if err != nil {
			return err
		}
		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		h.Name = rel
		h.Method = zip.Deflate
		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		if _, err := io.Copy(w, f); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func unzipToDir(data []byte, outDir string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, f := range r.File {
		dest := filepath.Join(outDir, filepath.FromSlash(f.Name))
		if !strings.HasPrefix(filepath.Clean(dest), filepath.Clean(outDir)+string(os.PathSeparator)) {
			return errors.New("zip path traversal detected")
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(dest)
		if err != nil {
			_ = rc.Close()
			return err
		}
		_, copyErr := io.Copy(out, rc)
		_ = out.Close()
		_ = rc.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func isZip(data []byte) bool {
	return len(data) >= 4 && bytes.Equal(data[:4], zipMagic)
}
