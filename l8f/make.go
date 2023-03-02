package l8f

import (
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func makeFramesLumpFile(inFramesDirectory, outFilePath string) (MakeVideoLumpTemp, error) {
	vlt := MakeVideoLumpTemp{}
	dirFIs, err := os.ReadDir(inFramesDirectory)
	if err != nil {
		return vlt, errors.Wrap(err, "os error")
	}

	inFrameNumbers := make([]int, 0)
	for _, dirFI := range dirFIs {
		if dirFI.IsDir() {
			return vlt, errors.New(fmt.Sprintf("the inFramesDirectory '%s' must not contain any subfolder", inFramesDirectory))
		}
		inFrameNameInt, err := strconv.Atoi(strings.ReplaceAll(dirFI.Name(), ".png", ""))
		if err != nil {
			return vlt, errors.New(fmt.Sprintf("the file '%s' of inFramesDirectory '%s' is not a number", dirFI.Name(), inFramesDirectory))
		}
		inFrameNumbers = append(inFrameNumbers, inFrameNameInt)
	}

	sort.Ints(inFrameNumbers)

	firstFrameHandle, err := os.Open(filepath.Join(inFramesDirectory, "1.png"))
	if err != nil {
		return vlt, errors.New(fmt.Sprintf("the inFramesDirectory '%s' has no '1.png'", inFramesDirectory))
	}
	im, _, err := image.DecodeConfig(firstFrameHandle)
	if err != nil {
		return vlt, errors.Wrap(err, "image error")
	}
	firstWidth := im.Width
	firstHeight := im.Height
	firstFrameHandle.Close()

	// validate same width and height of all the frames
	for _, inFrameNumber := range inFrameNumbers {
		currentFramePath := filepath.Join(inFramesDirectory, fmt.Sprintf("%d.png", inFrameNumber))
		currentFrameHandle, err := os.Open(currentFramePath)
		if err != nil {
			return vlt, errors.Wrap(err, "os error")
		}
		currentIm, _, err := image.DecodeConfig(currentFrameHandle)
		if err != nil {
			return vlt, errors.Wrap(err, "image error")
		}
		if currentIm.Width != firstWidth || currentIm.Height != firstHeight {
			return vlt, errors.New(fmt.Sprintf("the width or height of '%s' differs from the first frame", currentFramePath))
		}
	}

	// make temporary lump of unique frames
	outFileHandle, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return vlt, errors.Wrap(err, "os error")
	}
	uniqueFrames := make([]UniqueFrameDetails, 0) //first frame no and the size
	framesPointer := make(map[int]int)
	for _, inFrameNumber := range inFrameNumbers {
		currentFramePath := filepath.Join(inFramesDirectory, fmt.Sprintf("%d.png", inFrameNumber))
		currentFrameHandle, err := os.Open(currentFramePath)
		if err != nil {
			return vlt, errors.Wrap(err, "os error")
		}
		defer currentFrameHandle.Close()

		hashHandle := sha256.New()
		if _, err := io.Copy(hashHandle, currentFrameHandle); err != nil {
			return vlt, errors.Wrap(err, "io error")
		}
		hashOfCurrentFile := fmt.Sprintf("%x", hashHandle.Sum(nil))

		currentFrameHandle2, err := os.Open(currentFramePath)
		if err != nil {
			return vlt, errors.Wrap(err, "os error")
		}
		defer currentFrameHandle2.Close()

		ufq, err := findInUniqueFramesSlice(uniqueFrames, hashOfCurrentFile)
		if err == nil {
			framesPointer[inFrameNumber] = ufq.FirstFrameNumber
		} else {
			writtenSize, err := io.Copy(outFileHandle, currentFrameHandle2)
			if err != nil {
				return vlt, errors.Wrap(err, "io error")
			}
			uniqueFrames = append(uniqueFrames, UniqueFrameDetails{hashOfCurrentFile, inFrameNumber, int(writtenSize)})
			framesPointer[inFrameNumber] = inFrameNumber
		}
	}
	outFileHandle.Close()

	return MakeVideoLumpTemp{uniqueFrames, framesPointer}, nil
}

// MakeL8F is good for videos with a lot of stills eg. lyrics videos with a single background.
// the inFramesDirectory must contain png files numbered from 1.png upwards
// the framerate must be stored in the **meta** as a string
func MakeL8F(inFramesLaptopDirectory, inFramesMobileDirectory, inAudioFile string,
	meta map[string]string, tmpVideoDirectory, outFilePath string) error {
	if !doesPathExists(inFramesLaptopDirectory) {
		return errors.New(fmt.Sprintf("the path '%s' does not exists", inFramesLaptopDirectory))
	}
	if !strings.HasSuffix(inAudioFile, ".mp3") {
		return errors.New("The inAudioFile must be of type 'mp3'")
	}
	if !strings.HasSuffix(outFilePath, ".l8f") {
		return errors.New("The outFilePath must end with '.l8f'")
	}

	for k, v := range meta {
		if strings.Contains(k, "\n") || strings.Contains(v, "\n") {
			return errors.New("The meta elements must not contain newline")
		}
		if strings.Contains(k, ":") || strings.Contains(v, ":") {
			return errors.New("The meta elements must not contain ':' ")
		}
	}

	if _, ok := meta["framerate"]; !ok {
		return errors.New("the 'meta' map doesn't contain 'framerate'")
	}

	laptopLumpPath := filepath.Join(tmpVideoDirectory, ".tmp_"+untestedRandomString(10))
	mobileLumpPath := filepath.Join(tmpVideoDirectory, ".tmp_"+untestedRandomString(10))

	var wg sync.WaitGroup
	wg.Add(1)

	var lvlt MakeVideoLumpTemp
	var mvlt MakeVideoLumpTemp

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		out, err := makeFramesLumpFile(inFramesLaptopDirectory, laptopLumpPath)
		if err != nil {
			panic(err)
		}
		lvlt = out
	}(&wg)

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		out, err := makeFramesLumpFile(inFramesMobileDirectory, mobileLumpPath)
		if err != nil {
			panic(err)
		}
		mvlt = out
	}(&wg)

	wg.Wait()

	// write meta
	outStr := "meta:\n"
	for metaKey, metaValue := range meta {
		outStr += metaKey + ": " + metaValue + "\n"
	}
	outStr += "::\n"

	// write laptop_unique_frames
	outStr += "laptop_unique_frames:\n"
	for _, ufq := range lvlt.UniqueFrames {
		outStr += fmt.Sprintf("%d: %d\n", ufq.FirstFrameNumber, ufq.Size)
	}
	outStr += "::\n"

	// write laptop frames info
	outStr += "laptop_frames:\n"
	for frameNumber, pointedToFrameNumber := range lvlt.FramesPointerToUniqueFrames {
		outStr += fmt.Sprintf("%d: %d\n", frameNumber, pointedToFrameNumber)
	}
	outStr += "::\n"

	// write mobile_unique_frames
	outStr += "mobile_unique_frames:\n"
	for _, ufq := range mvlt.UniqueFrames {
		outStr += fmt.Sprintf("%d: %d\n", ufq.FirstFrameNumber, ufq.Size)
	}
	outStr += "::\n"

	// write mobile frames info
	outStr += "mobile_frames:\n"
	for frameNumber, pointedToFrameNumber := range mvlt.FramesPointerToUniqueFrames {
		outStr += fmt.Sprintf("%d: %d\n", frameNumber, pointedToFrameNumber)
	}
	outStr += "::\n"

	// write lumps
	outStr += "binary:\n"
	inAudioFileStat, err := os.Stat(inAudioFile)
	if err != nil {
		return errors.Wrap(err, "os error")
	}

	laptopLumpPathStat, err := os.Stat(laptopLumpPath)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	mobileLumpPathStat, err := os.Stat(mobileLumpPath)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	outStr += fmt.Sprintf("audio: %d\n", inAudioFileStat.Size())
	outStr += fmt.Sprintf("laptop_frames_lump: %d\n", laptopLumpPathStat.Size())
	outStr += fmt.Sprintf("mobile_frames_lump: %d\n", mobileLumpPathStat.Size())
	outStr += "::\n"

	outVideoHandle, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	defer outVideoHandle.Close()

	outVideoHandle.WriteString(fmt.Sprintf("%d\n", len(outStr)))
	outVideoHandle.WriteString(outStr)

	inAudioHandle, err := os.Open(inAudioFile)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	_, err = io.Copy(outVideoHandle, inAudioHandle)
	if err != nil {
		return errors.Wrap(err, "io error")
	}
	laptopLumpPathHandle, err := os.Open(laptopLumpPath)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	defer laptopLumpPathHandle.Close()
	_, err = io.Copy(outVideoHandle, laptopLumpPathHandle)
	if err != nil {
		return errors.Wrap(err, "io error")
	}

	mobileLumpPathHandle, err := os.Open(mobileLumpPath)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	defer mobileLumpPathHandle.Close()
	_, err = io.Copy(outVideoHandle, mobileLumpPathHandle)
	if err != nil {
		return errors.Wrap(err, "io error")
	}

	return nil
}

// the framerate must be stored in the **meta** as a string
func UpdateMeta(inVideoPath string, meta map[string]string, tmpVideoDirectory, outFilePath string) error {
	for k, v := range meta {
		if strings.Contains(k, "\n") || strings.Contains(v, "\n") {
			return errors.New("The meta elements must not contain newline")
		}
		if strings.Contains(k, ":") || strings.Contains(v, ":") {
			return errors.New("The meta elements must not contain ':' ")
		}
	}

	if !strings.HasSuffix(outFilePath, ".l8f") {
		return errors.New("The outFilePath must end with '.l8f'")
	}

	vhSize, err := getHeaderLengthFromVideo(inVideoPath)
	if err != nil {
		return err
	}
	vh, err := ReadHeaderFromVideo(inVideoPath)
	if err != nil {
		return err
	}

	vh.Meta = meta

	audioBytes := make([]byte, vh.AudioSize)

	rawVideoHandle, err := os.Open(inVideoPath)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	defer rawVideoHandle.Close()

	audioOffset := vhSize + 1 + len(fmt.Sprintf("%d", vhSize))
	_, err = rawVideoHandle.ReadAt(audioBytes, int64(audioOffset))
	if err != nil {
		return errors.Wrap(err, "strconv error")
	}

	laptopVideoBytes := make([]byte, vh.LaptopVideoSize)
	laptopVideoOffset := audioOffset + vh.AudioSize
	_, err = rawVideoHandle.ReadAt(laptopVideoBytes, int64(laptopVideoOffset))
	if err != nil {
		return errors.Wrap(err, "strconv error")
	}

	mobileVideoBytes := make([]byte, vh.MobileVideoSize)
	mobileVideoOffset := audioOffset + vh.AudioSize + vh.LaptopVideoSize
	_, err = rawVideoHandle.ReadAt(mobileVideoBytes, int64(mobileVideoOffset))
	if err != nil {
		return errors.Wrap(err, "strconv error")
	}

	// write meta
	outStr := "meta:\n"
	for metaKey, metaValue := range vh.Meta {
		outStr += metaKey + ": " + metaValue + "\n"
	}
	outStr += "::\n"

	// write unique_frames
	outStr += "laptop_unique_frames:\n"
	for _, ufq := range vh.LaptopUniqueFrames {
		outStr += fmt.Sprintf("%d: %d\n", ufq[0], ufq[1])
	}
	outStr += "::\n"

	// write frames info
	outStr += "laptop_frames:\n"
	for frameNumber, pointedToFrameNumber := range vh.LaptopFrames {
		outStr += fmt.Sprintf("%d: %d\n", frameNumber, pointedToFrameNumber)
	}
	outStr += "::\n"

	// write unique_frames
	outStr += "mobile_unique_frames:\n"
	for _, ufq := range vh.MobileUniqueFrames {
		outStr += fmt.Sprintf("%d: %d\n", ufq[0], ufq[1])
	}
	outStr += "::\n"

	// write frames info
	outStr += "mobile_frames:\n"
	for frameNumber, pointedToFrameNumber := range vh.MobileFrames {
		outStr += fmt.Sprintf("%d: %d\n", frameNumber, pointedToFrameNumber)
	}
	outStr += "::\n"

	// write lumps
	outStr += "binary:\n"
	outStr += fmt.Sprintf("audio: %d\n", vh.AudioSize)
	outStr += fmt.Sprintf("laptop_frames_lump: %d\n", vh.LaptopVideoSize)
	outStr += fmt.Sprintf("mobile_frames_lump: %d\n", vh.MobileVideoSize)
	outStr += "::\n"

	outVideoHandle, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	defer outVideoHandle.Close()

	outVideoHandle.WriteString(fmt.Sprintf("%d\n", len(outStr)))
	outVideoHandle.WriteString(outStr)

	outVideoHandle.Write(audioBytes)
	outVideoHandle.Write(laptopVideoBytes)
	outVideoHandle.Write(mobileVideoBytes)
	return nil
}
