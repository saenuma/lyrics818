package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/saenuma/lyrics818/l8shared"
)

//go:embed "ffplay.exe"
var FFPlayBytes []byte

func GetFFPlayCommand() string {
	homeDir, _ := os.UserHomeDir()

	ffmegDir := filepath.Join(homeDir, ".l818")
	outPath := filepath.Join(ffmegDir, "ffplay.exe")
	if !l8shared.DoesPathExists(outPath) {
		os.MkdirAll(ffmegDir, 0777)

		os.WriteFile(outPath, FFPlayBytes, 0777)
	}

	return outPath
}
