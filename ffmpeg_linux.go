package main

import (
	"os"
	"path/filepath"
)

func GetFFMPEGCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "ffmpeg"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", "ffmpeg")
	}

	return cmdPath
}
