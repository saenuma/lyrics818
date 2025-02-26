package main

import (
	"os"
	"path/filepath"

	internal "github.com/saenuma/lyrics818/internal/lyrics818"
)

func GetFFMPEGCommand() string {
	homeDir, _ := os.UserHomeDir()

	ffmegDir := filepath.Join(homeDir, ".l818")
	outPath := filepath.Join(ffmegDir, "ffmpeg.exe")
	if !internal.DoesPathExists(outPath) {
		os.MkdirAll(ffmegDir, 0777)

		os.WriteFile(outPath, ffmpegBytes, 0777)
	}

	return outPath
}
