package l8shared

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/saenuma/lyrics818/l8f"
)

func MakeVideo2(inputs map[string]string) (string, error) {
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

	laptopOutName := "lframes_" + time.Now().Format("20060102T150405")
	mobileOutName := "mframes_" + time.Now().Format("20060102T150405")
	lrenderPath := filepath.Join(rootPath, laptopOutName)
	os.MkdirAll(lrenderPath, 0777)
	mrenderPath := filepath.Join(rootPath, mobileOutName)
	os.MkdirAll(mrenderPath, 0777)

	MakeLaptopFrames(laptopOutName, totalSeconds, lrenderPath, inputs)
	MakeMobileFrames(mobileOutName, totalSeconds, mrenderPath, inputs)

	outName := strings.ReplaceAll(filepath.Base(fullMp3Path), ".mp3", ".l8f")
	fullOutPath := filepath.Join(rootPath, outName)
	err = l8f.MakeL8F(lrenderPath, mrenderPath, fullMp3Path, map[string]string{},
		rootPath, fullOutPath)
	if err != nil {
		panic(err)
	}
	os.RemoveAll(lrenderPath)
	os.RemoveAll(mrenderPath)

	return outName, nil
}
