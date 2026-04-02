package service

import (
	"bytes"
	"context"
	"fmt"
	"iziumov/tv-v-x/config"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DownloadResult struct {
	FilePath string
	FileSize int64
}

type DownloaderService struct {
	binaryPath  string
	outputDir   string
	maxFileSize string
	format      string
}

func NewDownloaderService(cfg config.YtdlpConfig) *DownloaderService {
	return &DownloaderService{
		binaryPath:  cfg.BinaryPath,
		outputDir:   cfg.OutputDir,
		maxFileSize: cfg.MaxFileSize,
		format:      cfg.Format,
	}
}

func (s *DownloaderService) Download(ctx context.Context, url string, jobID int) (*DownloadResult, error) {
	outputTemplate := filepath.Join(s.outputDir, fmt.Sprintf("job_%d.%%(ext)s", jobID))

	args := []string{
		"--no-playlist",
		"-f", s.format,
		"-o", outputTemplate,
		"--max-filesize", s.maxFileSize,
		"--no-overwrites",
		"--print", "after_move:filepath",
		url,
	}

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, s.binaryPath, args...)
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp error: %w, stderr: %s", err, stderr.String())
	}

	filePath := strings.TrimSpace(string(output))

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot stat donwloaded file: %w", err)
	}

	return &DownloadResult{
		FilePath: filePath,
		FileSize: info.Size(),
	}, nil
}

func (s *DownloaderService) CleanUp(filePath string) error {
	return os.Remove(filePath)
}
