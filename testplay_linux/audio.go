package main

import (
	"os"
	"path/filepath"
	"strings"
)

func GetMPCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "madplay"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "usr", "bin", "madplay")
	}

	return cmdPath
}
