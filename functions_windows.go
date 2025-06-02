package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sqweek/dialog"
)

func PickImageFile() string {
	filename, err := dialog.File().Filter("PNG Image", "png").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func PickTxtFile() string {
	filename, err := dialog.File().Filter("Lyrics File", "txt").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func PickFontFile() string {
	filename, err := dialog.File().Filter("Font file", "ttf").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func PickMp3File() string {
	filename, err := dialog.File().Filter("MP3 Audio", "mp3").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func GetFFMPEGCommand() string {
	execPath, _ := os.Executable()
	cmdPath := filepath.Join(filepath.Dir(execPath), "ffmpeg.exe")

	return cmdPath
}

func pickColor() string {
	execPath, _ := os.Executable()
	cmdPath := filepath.Join(filepath.Dir(execPath), "acpicker.exe")
	cmd := exec.Command(cmdPath)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}
