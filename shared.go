package main

import (
  "os"
  "strings"
  "github.com/pkg/errors"
  "path/filepath"
  "strconv"
  "fmt"
)



func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "os error")
	}
	dd := os.Getenv("SNAP_USER_COMMON")
	if strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) || dd == "" {
		dd = filepath.Join(hd, "lyrics818_data")
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


func parseLyricsFile(inPath string) map[int]string {
  raw, err := os.ReadFile(inPath)
  if err != nil {
    panic(err)
  }

  retObj := make(map[int]string)
  parts := strings.Split(string(raw), "\n\n")
  for _, part := range parts {
    innerParts := strings.Split(strings.TrimSpace(part), "\n")
    secs := timeFormatToSeconds(strings.TrimSpace(innerParts[0]))
    retObj[secs] = strings.Join(innerParts[1:], "\n")
  }

  return retObj
}


func getRenderPath(filename string) string {
	rootPath, _ := GetRootPath()
	added := 1
	for {
		f := filepath.Join(rootPath, fmt.Sprintf("%s_%d", filename, added))
		if DoesPathExists(f) {
			added += 1
		} else {
			os.MkdirAll(f, 0777)
			return f
		}
	}
}
