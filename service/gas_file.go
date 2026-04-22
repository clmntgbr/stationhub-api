package service

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"stationhub-api/config"
)

type GasFileService struct {
	config *config.Config
}

func NewGasFileService() *GasFileService {
	return &GasFileService{config: config.Load()}
}

func (s *GasFileService) DownloadGasFile() (string, error) {
	resp, err := http.Get(s.config.GasFileURL)
	if err != nil {
		return "", fmt.Errorf("failed to download the zip file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download the zip file: HTTP %d", resp.StatusCode)
	}

	zipFilePath := filepath.Join(os.TempDir(), "gas.zip")
	out, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save the zip file: %w", err)
	}

	return zipFilePath, nil
}

func (s *GasFileService) Extract(zipFilePath string) (string, error) {
	if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("ZIP file does not exist: %s", zipFilePath)
	}

	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer r.Close()

	if len(r.File) == 0 {
		return "", fmt.Errorf("ZIP file is empty")
	}

	firstFile := r.File[0]
	tempDir := os.TempDir()
	extractedPath := filepath.Join(tempDir, firstFile.Name)
	targetPath := filepath.Join(tempDir, "gas.xml")

	if err := s.extractFile(firstFile, extractedPath); err != nil {
		return "", fmt.Errorf("failed to extract file: %w", err)
	}

	// Remove existing gas.xml if present
	if _, err := os.Stat(targetPath); err == nil {
		if err := os.Remove(targetPath); err != nil {
			return "", fmt.Errorf("failed to remove existing gas.xml: %w", err)
		}
	}

	if err := os.Rename(extractedPath, targetPath); err != nil {
		return "", fmt.Errorf("failed to rename extracted file: %w", err)
	}

	return targetPath, nil
}

func (s *GasFileService) extractFile(f *zip.File, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}
