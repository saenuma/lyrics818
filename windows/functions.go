package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/saenuma/lyrics818/internal"
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

func pickColor() string {
	homeDir, _ := os.UserHomeDir()

	ffmegDir := filepath.Join(homeDir, ".l818")
	cmdPath := filepath.Join(ffmegDir, "acpicker.exe")
	if !internal.DoesPathExists(cmdPath) {
		os.MkdirAll(ffmegDir, 0777)

		os.WriteFile(cmdPath, acPickerBytes, 0777)
	}

	cmd := exec.Command(cmdPath)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}
