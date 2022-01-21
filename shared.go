package main


import (
  "os"
  "strings"
  "github.com/pkg/errors"
  "path/filepath"
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
