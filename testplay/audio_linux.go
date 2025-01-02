package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/saenuma/lyrics818/internal"
	"github.com/saenuma/lyrics818/l8f"
)

var playerCancelFn context.CancelFunc

func GetMPCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "madplay"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "usr", "bin", "madplay")
	}

	return cmdPath
}

func playAudio(l8fPath string) {
	rootPath, _ := internal.GetRootPath()
	mplayCmd := GetMPCommand()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	audioBytes, err := l8f.ReadAudio(l8fPath)
	if err != nil {
		panic(err)
	}

	tmpAudioPath := filepath.Join(rootPath, ".tmp_audio.mp3")
	os.WriteFile(tmpAudioPath, audioBytes, 0777)

	playerCancelFn = cancel
	exec.CommandContext(ctx, mplayCmd, tmpAudioPath).Run()
}
