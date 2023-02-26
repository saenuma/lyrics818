package l8f

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func getHeaderLengthFromVideo(inVideoPath string) (int, error) {
	if !doesPathExists(inVideoPath) {
		return 0, errors.New(fmt.Sprintf("the path '%s' does not exists", inVideoPath))
	}
	if !strings.HasSuffix(inVideoPath, ".l8f") {
		return 0, errors.New("The inVideoPath must be of type 'l8f'")
	}

	rawVideoHandle, err := os.Open(inVideoPath)
	if err != nil {
		return 0, errors.Wrap(err, "os error")
	}
	defer rawVideoHandle.Close()

	var count int64
	var headerLengthStr string
	for {
		inByte := make([]byte, 1)
		_, err := rawVideoHandle.ReadAt(inByte, count)
		if err != nil {
			return 0, errors.Wrap(err, "os error")
		}
		if string(inByte) != "\n" {
			headerLengthStr += string(inByte)
			count += 1
		} else {
			break
		}
		continue
	}

	headerLength, err := strconv.Atoi(headerLengthStr)
	if err != nil {
		return 0, errors.Wrap(err, "strconv error")
	}

	return headerLength, nil
}

func ReadHeaderFromVideo(inVideoPath string) (VideoHeader, error) {
	evh := VideoHeader{}
	if !doesPathExists(inVideoPath) {
		return evh, errors.New(fmt.Sprintf("the path '%s' does not exists", inVideoPath))
	}
	if !strings.HasSuffix(inVideoPath, ".l8f") {
		return evh, errors.New("The inVideoPath must be of type 'l8f'")
	}

	rawVideoHandle, err := os.Open(inVideoPath)
	if err != nil {
		return evh, errors.Wrap(err, "os error")
	}
	defer rawVideoHandle.Close()

	headerLength, err := getHeaderLengthFromVideo(inVideoPath)
	if err != nil {
		return evh, errors.Wrap(err, "strconv error")
	}
	headerBytes := make([]byte, headerLength)
	readBegin := int64(len(strconv.Itoa(headerLength))) + 1
	_, err = rawVideoHandle.ReadAt(headerBytes, readBegin)
	if err != nil {
		return evh, errors.Wrap(err, "os error")
	}
	headerStr := string(headerBytes)

	// begin parsing Video header
	metaBeginPart := strings.Index(headerStr, "meta:")
	metaEndPart := strings.Index(headerStr[metaBeginPart:], "::")
	if metaEndPart == -1 {
		return evh, errors.New("Bad Header: meta section must end with a '::'.")
	}
	metaPart := headerStr[metaBeginPart+len("meta:\n") : metaBeginPart+metaEndPart]
	meta := make(map[string]string)
	for _, line := range strings.Split(metaPart, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		partsOfLine := strings.Split(line, ":")
		meta[partsOfLine[0]] = strings.TrimSpace(partsOfLine[1])
	}
	evh.Meta = meta

	luniqueFramesBeginPart := strings.Index(headerStr, "laptop_unique_frames:")
	luniqueFramesEndPart := strings.Index(headerStr[luniqueFramesBeginPart:], "::")
	if luniqueFramesEndPart == -1 {
		return evh, errors.New("Bad Header: laptop_unique_frames section must end with a '::'.")
	}
	luniqueFramesPart := headerStr[luniqueFramesBeginPart+len("laptop_unique_frames:\n") : luniqueFramesBeginPart+luniqueFramesEndPart]
	luniqueFrames := make([][]int, 0)
	for _, line := range strings.Split(luniqueFramesPart, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		partsOfLine := strings.Split(line, ":")
		f1, err := strconv.Atoi(strings.TrimSpace(partsOfLine[0]))
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}
		f2, err := strconv.Atoi(strings.TrimSpace(partsOfLine[1]))
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}

		luniqueFrames = append(luniqueFrames, []int{f1, f2})
	}
	evh.LaptopUniqueFrames = luniqueFrames

	lframesBeginPart := strings.LastIndex(headerStr, "laptop_frames:")
	lframesEndPart := strings.Index(headerStr[lframesBeginPart:], "::")
	if lframesEndPart == -1 {
		return evh, errors.New("Bad Header: laptop_frames section must end with a '::'.")
	}
	lframesPart := headerStr[lframesBeginPart+len("laptop_frames:\n") : lframesEndPart+lframesBeginPart]
	lframes := make(map[int]int)
	for _, line := range strings.Split(lframesPart, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		partsOfLine := strings.Split(line, ":")
		frame1Int, err := strconv.Atoi(partsOfLine[0])
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}
		frame2Int, err := strconv.Atoi(strings.TrimSpace(partsOfLine[1]))
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}
		lframes[frame1Int] = frame2Int
	}
	evh.LaptopFrames = lframes

	muniqueFramesBeginPart := strings.Index(headerStr, "mobile_unique_frames:")
	muniqueFramesEndPart := strings.Index(headerStr[muniqueFramesBeginPart:], "::")
	if muniqueFramesEndPart == -1 {
		return evh, errors.New("Bad Header: mobile_unique_frames section must end with a '::'.")
	}
	muniqueFramesPart := headerStr[muniqueFramesBeginPart+len("mobile_unique_frames:\n") : muniqueFramesBeginPart+muniqueFramesEndPart]
	muniqueFrames := make([][]int, 0)
	for _, line := range strings.Split(muniqueFramesPart, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		partsOfLine := strings.Split(line, ":")
		f1, err := strconv.Atoi(strings.TrimSpace(partsOfLine[0]))
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}
		f2, err := strconv.Atoi(strings.TrimSpace(partsOfLine[1]))
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}

		muniqueFrames = append(muniqueFrames, []int{f1, f2})
	}
	evh.MobileUniqueFrames = muniqueFrames

	mframesBeginPart := strings.LastIndex(headerStr, "mobile_frames:")
	mframesEndPart := strings.Index(headerStr[mframesBeginPart:], "::")
	if mframesEndPart == -1 {
		return evh, errors.New("Bad Header: mobile_frames section must end with a '::'.")
	}
	mframesPart := headerStr[mframesBeginPart+len("mobile_frames:\n") : mframesEndPart+mframesBeginPart]
	mframes := make(map[int]int)
	for _, line := range strings.Split(mframesPart, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		partsOfLine := strings.Split(line, ":")
		frame1Int, err := strconv.Atoi(partsOfLine[0])
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}
		frame2Int, err := strconv.Atoi(strings.TrimSpace(partsOfLine[1]))
		if err != nil {
			return evh, errors.Wrap(err, "strconv error")
		}
		mframes[frame1Int] = frame2Int
	}
	evh.MobileFrames = mframes

	binaryBeginPart := strings.Index(headerStr, "binary:")
	binaryEndPart := strings.Index(headerStr[binaryBeginPart:], "::")
	if binaryBeginPart == -1 {
		return evh, errors.New("Bad Header: chapters section be present.")
	}
	binaryPart := headerStr[binaryBeginPart+len("binary:\n") : binaryBeginPart+binaryEndPart]
	partsOfBinaryPart := strings.Split(binaryPart, "\n")
	audioPart := partsOfBinaryPart[0]
	lvideoPart := partsOfBinaryPart[1]
	mVideoPart := partsOfBinaryPart[2]

	audioSizeStr := audioPart[len("audio: "):]
	audioSizeInt, err := strconv.Atoi(audioSizeStr)
	if err != nil {
		return evh, errors.Wrap(err, "strconv error")
	}
	lvideoSizeStr := lvideoPart[len("laptop_frames_lump: "):]
	lvideoSizeInt, err := strconv.Atoi(lvideoSizeStr)
	if err != nil {
		return evh, errors.Wrap(err, "strconv error")
	}
	mvideoSizeStr := mVideoPart[len("mobile_frames_lump: "):]
	mvideoSizeInt, err := strconv.Atoi(mvideoSizeStr)
	if err != nil {
		return evh, errors.Wrap(err, "strconv error")
	}

	evh.AudioSize = audioSizeInt
	evh.LaptopVideoSize = lvideoSizeInt
	evh.MobileVideoSize = mvideoSizeInt

	return evh, nil
}

// The audio is []bytes but it should contain 'mp3' audio
func ReadAudio(inVideoPath string) ([]byte, error) {
	vhSize, err := getHeaderLengthFromVideo(inVideoPath)
	if err != nil {
		return nil, err
	}
	vh, err := ReadHeaderFromVideo(inVideoPath)
	if err != nil {
		return nil, err
	}

	audioBytes := make([]byte, vh.AudioSize)

	rawVideoHandle, err := os.Open(inVideoPath)
	if err != nil {
		return nil, errors.Wrap(err, "os error")
	}
	defer rawVideoHandle.Close()

	audioOffset := vhSize + 1 + len(fmt.Sprintf("%d", vhSize))
	_, err = rawVideoHandle.ReadAt(audioBytes, int64(audioOffset))
	if err != nil {
		return nil, errors.Wrap(err, "strconv error")
	}

	return audioBytes, nil
}

// Read frames for 1 seconds, starting from the 'seconds' parameter
// 'seconds' parameter starts from 0
func ReadLaptopFrames(inVideoPath string, seconds int) ([]*image.Image, error) {
	vhSize, err := getHeaderLengthFromVideo(inVideoPath)
	if err != nil {
		return nil, err
	}

	vh, err := ReadHeaderFromVideo(inVideoPath)
	if err != nil {
		return nil, err
	}

	rawVideoHandle, err := os.Open(inVideoPath)
	if err != nil {
		return nil, errors.Wrap(err, "os error")
	}
	defer rawVideoHandle.Close()

	audioOffset := vhSize + 1 + len(fmt.Sprintf("%d", vhSize))
	videoBytes := make([]byte, vh.LaptopVideoSize)
	videoOffset := audioOffset + vh.AudioSize
	_, err = rawVideoHandle.ReadAt(videoBytes, int64(videoOffset))
	if err != nil {
		return nil, errors.Wrap(err, "strconv error")
	}

	allFrames := make([]int, 0)
	for k := range vh.LaptopFrames {
		allFrames = append(allFrames, k)
	}

	sort.Ints(allFrames)

	frameRate := vh.Meta["framerate"]
	frameRateInt, err := strconv.Atoi(frameRate)
	if err != nil {
		return nil, errors.Wrap(err, "strconv error")
	}
	toRetFrames := allFrames[seconds*frameRateInt : (seconds+1)*frameRateInt]

	images := make([]*image.Image, 0)
	for _, frameNumber := range toRetFrames {
		pointedToFrameNumber := vh.LaptopFrames[frameNumber]

		// unpack the right frame
		readFrameOffset := 0
		toReadSize := 0

		for _, uniqueFrameDetails := range vh.LaptopUniqueFrames {
			if uniqueFrameDetails[0] == pointedToFrameNumber {
				toReadSize = int(uniqueFrameDetails[1])
				break
			} else {
				readFrameOffset += int(uniqueFrameDetails[1])
			}
		}

		currentFrameBytes := videoBytes[readFrameOffset : readFrameOffset+toReadSize]
		img, _, err := image.Decode(bytes.NewReader(currentFrameBytes))
		if err != nil {
			return nil, errors.Wrap(err, "image error")
		}

		images = append(images, &img)
	}

	return images, nil
}

// Read frames for 1 seconds, starting from the 'seconds' parameter
// 'seconds' parameter starts from 0
func ReadMobileFrames(inVideoPath string, seconds int) ([]*image.Image, error) {
	vhSize, err := getHeaderLengthFromVideo(inVideoPath)
	if err != nil {
		return nil, err
	}

	vh, err := ReadHeaderFromVideo(inVideoPath)
	if err != nil {
		return nil, err
	}

	rawVideoHandle, err := os.Open(inVideoPath)
	if err != nil {
		return nil, errors.Wrap(err, "os error")
	}
	defer rawVideoHandle.Close()

	audioOffset := vhSize + 1 + len(fmt.Sprintf("%d", vhSize))
	videoBytes := make([]byte, vh.MobileVideoSize)
	videoOffset := audioOffset + vh.AudioSize + vh.LaptopVideoSize
	_, err = rawVideoHandle.ReadAt(videoBytes, int64(videoOffset))
	if err != nil {
		return nil, errors.Wrap(err, "strconv error")
	}

	allFrames := make([]int, 0)
	for k := range vh.MobileFrames {
		allFrames = append(allFrames, k)
	}

	sort.Ints(allFrames)

	frameRate := vh.Meta["framerate"]
	frameRateInt, err := strconv.Atoi(frameRate)
	if err != nil {
		return nil, errors.Wrap(err, "strconv error")
	}
	toRetFrames := allFrames[seconds*frameRateInt : (seconds+1)*frameRateInt]

	images := make([]*image.Image, 0)
	for _, frameNumber := range toRetFrames {
		pointedToFrameNumber := vh.MobileFrames[frameNumber]

		// unpack the right frame
		readFrameOffset := 0
		toReadSize := 0

		for _, uniqueFrameDetails := range vh.MobileUniqueFrames {
			if uniqueFrameDetails[0] == pointedToFrameNumber {
				toReadSize = int(uniqueFrameDetails[1])
				break
			} else {
				readFrameOffset += int(uniqueFrameDetails[1])
			}
		}

		currentFrameBytes := videoBytes[readFrameOffset : readFrameOffset+toReadSize]
		img, _, err := image.Decode(bytes.NewReader(currentFrameBytes))
		if err != nil {
			return nil, errors.Wrap(err, "image error")
		}

		images = append(images, &img)
	}

	return images, nil
}

// This checks the length of the video using the frames itself
// It doesn't check against the audio data embedded in it
func GetVideoLength(inVideoPath string) (int, error) {
	if !doesPathExists(inVideoPath) {
		return 0, errors.New(fmt.Sprintf("the path '%s' does not exists", inVideoPath))
	}
	if !strings.HasSuffix(inVideoPath, ".l8f") {
		return 0, errors.New("The inVideoPath must be of type 'l8f'")
	}

	vh, err := ReadHeaderFromVideo(inVideoPath)
	if err != nil {
		return 0, err
	}

	frameRate := vh.Meta["framerate"]
	frameRateInt, err := strconv.Atoi(frameRate)
	if err != nil {
		return 0, errors.Wrap(err, "strconv error")
	}

	totalSeconds := float64(len(vh.LaptopFrames)) / float64(frameRateInt)

	return int(math.Ceil(totalSeconds)), nil
}
