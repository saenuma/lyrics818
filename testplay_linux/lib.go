package main

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFFPlayCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "ffplay"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", "ffplay")
	}

	return cmdPath
}
