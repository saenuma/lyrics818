package l8shared

import (
	"fmt"
	"math"
	"os"

	"path/filepath"
	"strconv"
	"strings"

	"io"

	"github.com/pkg/errors"

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

func SecondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	secondsStr := fmt.Sprintf("%d", seconds)
	if seconds < 10 {
		secondsStr = "0" + secondsStr
	}
	str := fmt.Sprintf("%d:%s", minutes, secondsStr)
	return str
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

func FindIn(container []int, elem int) int {
	for i, o := range container {
		if o > elem {
			return i
		}
	}
	return -1
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
