package main

import (
	"os"
	// "fmt"
	"image"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/saenuma/zazabul"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "os error")
	}

	dd := os.Getenv("SNAP_USER_COMMON")

	if strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) || dd == "" {
		dd = filepath.Join(hd, "Lyrics818")
		os.MkdirAll(dd, 0777)
	}

	return dd, nil
}

func timeFormatToSeconds(s string) int {
	// calculate total duration of the song
	parts := strings.Split(s, ":")
	minutesPartConverted, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	secondsPartConverted, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	totalSecondsOfSong := (60 * minutesPartConverted) + secondsPartConverted
	return totalSecondsOfSong
}

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func parseLyricsFile(inPath string, totalSeconds int) map[int]string {
	raw, err := os.ReadFile(inPath)
	if err != nil {
		panic(err)
	}

	tmpObj := make(map[int]string)
	parts := strings.Split(string(raw), "\r\n\r\n")
	for _, part := range parts {
		innerParts := strings.Split(strings.TrimSpace(part), "\r\n")
		secs := timeFormatToSeconds(strings.TrimSpace(innerParts[0]))
		tmpObj[secs] = strings.Join(innerParts[1:], "\r\n")
	}

	retObj := make(map[int]string)
	started := false
	var lastSecondsWithLyrics int
	for seconds := 0; seconds < totalSeconds; seconds++ {
		if !started {
			txt, ok := tmpObj[seconds]
			if !ok {
				retObj[seconds] = ""
			} else {
				started = true
				retObj[seconds] = txt
				lastSecondsWithLyrics = seconds
			}

		} else {
			txt, ok := tmpObj[seconds]
			if !ok {
				retObj[seconds] = tmpObj[lastSecondsWithLyrics]
			} else {
				retObj[seconds] = txt
				lastSecondsWithLyrics = seconds
			}
		}
	}
	return retObj
}

func wordWrap(conf zazabul.Config, text string, writeWidth int) []string {
	rootPath, _ := GetRootPath()

	rgba := image.NewRGBA(image.Rect(0, 0, 1366, 768))

	fontBytes, err := os.ReadFile(filepath.Join(rootPath, conf.Get("font_file")))
	if err != nil {
		panic(err)
	}
	fontParsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	fontDrawer := &font.Drawer{
		Dst: rgba,
		Src: image.Black,
		Face: truetype.NewFace(fontParsed, &truetype.Options{
			Size:    SIZE,
			DPI:     DPI,
			Hinting: font.HintingNone,
		}),
	}

	widthFixed := fixed.I(writeWidth)

	strs := strings.Fields(text)
	outStrs := make([]string, 0)
	var tmpStr string
	for i, oneStr := range strs {
		var aStr string
		if i == 0 {
			aStr = oneStr
		} else {
			aStr += " " + oneStr
		}

		tmpStr += aStr
		if fontDrawer.MeasureString(tmpStr) >= widthFixed {
			outStr := tmpStr[:len(tmpStr)-len(aStr)]
			tmpStr = oneStr
			outStrs = append(outStrs, outStr)
		}
	}
	outStrs = append(outStrs, tmpStr)

	return outStrs
}

func FindIn(container []int, elem int) int {
	for i, o := range container {
		if o > elem {
			return i
		}
	}
	return -1
}

func GetFFMPEGCommand() string {
	// get the right ffmpeg command
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	devPath := filepath.Join(homeDir, "bin", "ffmpeg.exe")
	bundledPath := filepath.Join("C:\\Program Files (x86)\\Lyrics818", "ffmpeg.exe")
	if DoesPathExists(devPath) {
		return devPath
	}

	return bundledPath
}
