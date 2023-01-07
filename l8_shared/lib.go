package l8_shared

import (
	"fmt"
	"math"
	"os"

	"image"
	"path/filepath"
	"strconv"
	"strings"

	"io"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/saenuma/zazabul"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/tcolgate/mp3"
)

const (
	DPI     = 72.0
	SIZE    = 80.0
	SPACING = 1.1
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

func TimeFormatToSeconds(s string) int {
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

func ParseLyricsFile(inPath string, totalSeconds int) map[int]string {
	raw, err := os.ReadFile(inPath)
	if err != nil {
		panic(err)
	}

	tmpObj := make(map[int]string)
	cleanedLyricsStr := strings.ReplaceAll(string(raw), "\r\n", "\n")
	parts := strings.Split(cleanedLyricsStr, "\n\n")
	for _, part := range parts {
		innerParts := strings.Split(strings.TrimSpace(part), "\n")
		secs := TimeFormatToSeconds(strings.TrimSpace(innerParts[0]))
		tmpObj[secs] = strings.Join(innerParts[1:], "\n")
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

func ValidateLyrics(conf zazabul.Config, lyricsObject map[int]string) error {
	totalSeconds := TimeFormatToSeconds(conf.Get("total_length"))

	// validate the length of a page of lyrics
	for i := 1; i < totalSeconds; i++ {
		text := lyricsObject[i]
		texts := strings.Split(text, "\n")

		finalTexts := make([]string, 0)
		for _, txt := range texts {
			wrappedTxts := WordWrap(conf, txt, 1366-130)
			finalTexts = append(finalTexts, wrappedTxts...)
		}

		if len(finalTexts) > 7 {
			return errors.New(fmt.Sprintf("Shorten the following text for it to fit this video:\n%s",
				text))
		}
	}

	return nil
}

func WordWrap(conf zazabul.Config, text string, writeWidth int) []string {
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
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "ffmpeg"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", "ffmpeg")
	}

	return cmdPath
}

func ReadSecondsFromMusicFile(musicPath string) (int, error) {
	t := 0.0

	r, err := os.Open(musicPath)
	if err != nil {
		return 0, err
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		t = t + f.Duration().Seconds()
	}

	correctedT := math.Ceil(t)
	return int(correctedT), nil
}
