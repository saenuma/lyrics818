package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/saenuma/lyrics818/internal"
	"github.com/saenuma/lyrics818/l8f"
)

func playAudio(l8fPath string) {
	rootPath, _ := internal.GetRootPath()

	audioBytes, err := l8f.ReadAudio(l8fPath)
	if err != nil {
		panic(err)
	}

	tmpAudioPath := filepath.Join(rootPath, ".tmp_audio.mp3")
	os.WriteFile(tmpAudioPath, audioBytes, 0777)

	mpg := GetMPGCommand()

	ctx, cancel := context.WithCancel(context.Background())
	linuxCancelFn = cancel
	out, err := exec.CommandContext(ctx, mpg, tmpAudioPath).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
	}
}

func GetMPGCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "madplay"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "usr", "bin", "madplay")
	}

	return cmdPath
}
