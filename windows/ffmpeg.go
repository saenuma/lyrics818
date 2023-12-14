package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/saenuma/lyrics818/l8shared"
)

//go:embed "ffmpeg/ffmpeg.exe"
var ffmpegBytes []byte

func GetFFMPEGCommand() string {
	homeDir, _ := os.UserHomeDir()

	ffmegDir := filepath.Join(homeDir, ".l818")
	outPath := filepath.Join(ffmegDir, "ffmpeg.exe")
	if !l8shared.DoesPathExists(outPath) {
		os.MkdirAll(ffmegDir, 0777)

		os.WriteFile(outPath, ffmpegBytes, 0777)
	}

	return outPath
}
