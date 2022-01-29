package main

import (
  "os"
  "strings"
  "github.com/pkg/errors"
  "path/filepath"
  "strconv"
  "github.com/golang/freetype"
  "golang.org/x/image/font"
  "github.com/golang/freetype/truetype"
  "golang.org/x/image/math/fixed"
  "github.com/bankole7782/zazabul"
  "image"
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


func parseLyricsFile(inPath string, totalSeconds int) map[int]string {
  raw, err := os.ReadFile(inPath)
  if err != nil {
    panic(err)
  }

  tmpObj := make(map[int]string)
  parts := strings.Split(string(raw), "\n\n")
  for _, part := range parts {
    innerParts := strings.Split(strings.TrimSpace(part), "\n")
    secs := timeFormatToSeconds(strings.TrimSpace(innerParts[0]))
    tmpObj[secs] = strings.Join(innerParts[1:], "\n")
  }

  retObj := make(map[int]string)
  started := false
  var lastSecondsWithLyrics int
  for seconds := 0; seconds < totalSeconds; seconds++ {
    if started == false {
      txt, ok := tmpObj[seconds]
      if ! ok {
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
      Size: SIZE,
      DPI: DPI,
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
      outStr := tmpStr[ : len(tmpStr) - len(aStr) ]
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
