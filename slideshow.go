package main

import (
	"image"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/disintegration/imaging"
	color2 "github.com/gookit/color"
	"github.com/saenuma/zazabul"
)

// immediate appearance slideshow method
func MakeSlideshowFrames(outName string, totalSeconds int, renderPath string, conf zazabul.Config) {
	rootPath, _ := GetRootPath()

	fullPicsPath := filepath.Join(rootPath, conf.Get("pictures_dir"))
	if !DoesPathExists(fullPicsPath) {
		color2.Red.Printf("The pictures folder '%s' does not exist.\n Exiting.\n", fullPicsPath)
		os.Exit(1)
	}

	dirFIs, err := os.ReadDir(fullPicsPath)
	if err != nil {
		color2.Red.Printf("Error reading directory '%s'.\nFull Error: %s\n.Exiting", fullPicsPath, err.Error())
		os.Exit(1)
	}
	picsPaths := make([]string, 0)
	picsBytes := make(map[int]image.Image)
	for i, dirFI := range dirFIs {
		aPicPath := filepath.Join(fullPicsPath, dirFI.Name())
		aPicOpened, _ := imaging.Open(aPicPath)
		if aPicOpened.Bounds().Dx() != 1366 || aPicOpened.Bounds().Dy() != 768 {
			color2.Red.Printf("The width of the picture '%s'\n is not 1366px or the height is not 768px.\nExiting.\n", aPicPath)
			os.Exit(1)
		}
		picsBytes[i] = aPicOpened
		picsPaths = append(picsPaths, aPicPath)
	}

	// var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	var wg sync.WaitGroup

	switchFrequency := 15 // seconds
	totalThreadsF64 := float64(totalSeconds) / float64(switchFrequency)
	totalThreads := int(math.Floor(totalThreadsF64))

	for threadIndex := 0; threadIndex < totalThreads; threadIndex++ {
		wg.Add(1)

		startSeconds := threadIndex * switchFrequency
		endSeconds := (threadIndex + 1) * switchFrequency

		lengthOfPics := len(picsPaths)
		currentIndexF64 := math.Mod(float64(threadIndex), float64(lengthOfPics))
		currentIndex := int(currentIndexF64)

		go func(startSeconds, endSeconds, currentIndex int, wg *sync.WaitGroup) {
			defer wg.Done()

			for seconds := startSeconds; seconds < endSeconds; seconds++ {
				for i := 1; i <= 24; i++ {
					out := (24 * seconds) + i
					outPath := filepath.Join(renderPath, strconv.Itoa(out)+".png")

					imaging.Save(picsBytes[currentIndex], outPath)
				}
			}

		}(startSeconds, endSeconds, currentIndex, &wg)
	}
	wg.Wait()

	for seconds := (totalThreads * switchFrequency); seconds < totalSeconds; seconds++ {
		lengthOfPics := len(picsPaths)
		currentIndexF64 := math.Mod(float64(1+(totalThreads*switchFrequency)), float64(lengthOfPics))
		currentIndex := int(currentIndexF64)

		for i := 1; i <= 24; i++ {
			out := (24 * seconds) + i
			outPath := filepath.Join(renderPath, strconv.Itoa(out)+".png")

			imaging.Save(picsBytes[currentIndex], outPath)
		}

	}

}
