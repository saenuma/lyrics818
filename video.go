package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func MakeVideo(inputs map[string]string, ffmpegCommandPath string) (string, error) {

	rootPath, err := GetRootPath()
	if err != nil {
		return "", err
	}

	fullMp3Path := inputs["music_file"]
	if !strings.HasSuffix(fullMp3Path, ".mp3") {
		return "", err
	}

	totalSeconds, err := ReadSecondsFromMusicFile(fullMp3Path)
	if err != nil {
		return "", err
	}

	outName := ".frames_" + time.Now().Format("20060102T150405")

	renderPath := filepath.Join(rootPath, outName)
	os.MkdirAll(renderPath, 0777)

	// command := GetFFMPEGCommand()

	MakeLaptopFrames(outName, totalSeconds, renderPath, inputs)

	// make video from laptop frames
	_, err = exec.Command(ffmpegCommandPath, "-framerate", "1", "-i", filepath.Join(renderPath, "%d.png"),
		"-pix_fmt", "yuv420p",
		filepath.Join(renderPath, "tmp_"+outName+".mp4")).CombinedOutput()
	if err != nil {
		return "", err
	}

	videoFileName := strings.ReplaceAll(filepath.Base(fullMp3Path), ".mp3", ".mp4")
	fullOutPath := filepath.Join(rootPath, videoFileName)
	// join audio to video
	_, err = exec.Command(ffmpegCommandPath, "-y", "-i", filepath.Join(renderPath, "tmp_"+outName+".mp4"),
		"-i", inputs["music_file"], "-pix_fmt", "yuv420p", fullOutPath).CombinedOutput()
	if err != nil {
		return "", err
	}

	os.RemoveAll(renderPath)
	return fullOutPath, nil
}
