package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetExecPath(execName string) string {
	homeDir, _ := os.UserHomeDir()
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = filepath.Join(homeDir, "bin", execName)
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", execName)
	}

	return cmdPath
}

func pickFile(exts string) string {
	fPickerPath := GetExecPath("fpicker")

	rootPath, _ := GetRootPath()
	cmd := exec.Command(fPickerPath, rootPath, exts)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func PickImageFile() string {
	return pickFile("png")
}

func PickTxtFile() string {
	return pickFile("txt")
}

func PickFontFile() string {
	return pickFile("ttf")
}

func PickMp3File() string {
	return pickFile("mp3")
}

func pickColor() string {
	homeDir, _ := os.UserHomeDir()
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = filepath.Join(homeDir, "bin", "acpicker")
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", "acpicker")
	}

	cmd := exec.Command(cmdPath)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetFFMPEGCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "ffmpeg"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", "ffmpeg")
	}

	return cmdPath
}
