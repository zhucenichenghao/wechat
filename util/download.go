package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download 下载文件到临时目录
func Download(ctx context.Context, URL string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file status: %s", resp.Status)
	}

	tmpDir := os.TempDir()

	fileName := filepath.Base(URL)
	tmpFilePath := filepath.Join(tmpDir, fileName)

	out, err := os.Create(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create tmp file")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write tmp file")
	}

	return tmpFilePath, nil
}
