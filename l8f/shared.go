package l8f

import (
	"math/rand"
	"os"

	"github.com/pkg/errors"
)

func doesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func untestedRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz1234567890"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func findInUniqueFramesSlice(container []UniqueFrameDetails, hash string) (UniqueFrameDetails, error) {
	for _, ufq := range container {
		if hash == ufq.Hash {
			return ufq, nil
		}
	}

	return UniqueFrameDetails{}, errors.New("frame not found")
}
