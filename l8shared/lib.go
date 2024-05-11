package l8shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/pkg/errors"
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

type LongEntry struct{}

func (d *LongEntry) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 30)
}

func (d *LongEntry) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	// for _, o := range objects {
	newHeight := containerSize.Height
	newSize := fyne.NewSize(containerSize.Width, newHeight)
	objects[0].Resize(newSize)
	objects[0].Move(pos)
	// pos = pos.Add(fyne.NewPos(0, newHeight+10))
	// }
}

func GetFilesOfType(rootPath, ext string) []string {
	dirFIs, err := os.ReadDir(rootPath)
	if err != nil {
		panic(err)
	}
	files := make([]string, 0)
	for _, dirFI := range dirFIs {
		if !dirFI.IsDir() && !strings.HasPrefix(dirFI.Name(), ".") && strings.HasSuffix(dirFI.Name(), ext) {
			files = append(files, dirFI.Name())
		}

		if dirFI.IsDir() && !strings.HasPrefix(dirFI.Name(), ".") {
			innerDirFIs, _ := os.ReadDir(filepath.Join(rootPath, dirFI.Name()))
			innerFiles := make([]string, 0)

			for _, innerDirFI := range innerDirFIs {
				if !innerDirFI.IsDir() && !strings.HasPrefix(innerDirFI.Name(), ".") && strings.HasSuffix(innerDirFI.Name(), ext) {
					innerFiles = append(innerFiles, filepath.Join(dirFI.Name(), innerDirFI.Name()))
				}
			}

			if len(innerFiles) > 0 {
				files = append(files, innerFiles...)
			}
		}

	}

	return files
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
