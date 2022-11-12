package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	color2 "github.com/gookit/color"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/saenuma/zazabul"
	"golang.org/x/image/font"
)

func ImageMethod(outName string, totalSeconds int, renderPath string, conf zazabul.Config) {
	numberOfCPUS := runtime.NumCPU()
	rootPath, _ := GetRootPath()
	lyricsObject := ParseLyricsFile(filepath.Join(rootPath, conf.Get("lyrics_file")), totalSeconds)

	err := ValidateLyrics(conf, lyricsObject)
	if err != nil {
		color2.Red.Println(err)
		os.Exit(1)
	}

	jobsPerThread := int(math.Floor(float64(totalSeconds) / float64(numberOfCPUS)))
	// remainder := int(math.Mod(float64(totalSeconds), float64(numberOfCPUS)))
	var wg sync.WaitGroup

	for threadIndex := 0; threadIndex < numberOfCPUS; threadIndex++ {
		wg.Add(1)

		startSeconds := threadIndex * jobsPerThread
		endSeconds := (threadIndex + 1) * jobsPerThread

		go func(startSeconds, endSeconds int, wg *sync.WaitGroup) {
			defer wg.Done()

			for seconds := startSeconds; seconds < endSeconds; seconds++ {
				txt := lyricsObject[seconds]
				if txt == "" {
					img, err := imaging.Open(filepath.Join(rootPath, conf.Get("background_file")))
					if err != nil {
						panic(err)
					}
					writeManyImagesToDisk(img, renderPath, seconds)
				} else {
					img := writeLyricsToImage(conf, lyricsObject[seconds])
					writeManyImagesToDisk(img, renderPath, seconds)
				}

			}

		}(startSeconds, endSeconds, &wg)
	}
	wg.Wait()

	for seconds := (jobsPerThread * numberOfCPUS); seconds < totalSeconds; seconds++ {
		txt := lyricsObject[seconds]
		if txt == "" {
			img, err := imaging.Open(filepath.Join(rootPath, conf.Get("background_file")))
			if err != nil {
				panic(err)
			}
			writeManyImagesToDisk(img, renderPath, seconds)
		} else {
			img := writeLyricsToImage(conf, lyricsObject[seconds])
			writeManyImagesToDisk(img, renderPath, seconds)
		}
	}

	color2.Green.Println("Completed generating frames of your lyrics video")

	command := GetFFMPEGCommand()

	out, err := exec.Command(command, "-framerate", "24", "-i", filepath.Join(renderPath, "%d.png"),
		"-pix_fmt", "yuv420p",
		filepath.Join(renderPath, "tmp_"+outName+".mp4")).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

}

func writeManyImagesToDisk(img image.Image, renderPath string, seconds int) {
	for i := 1; i <= 24; i++ {
		out := (24 * seconds) + i
		outPath := filepath.Join(renderPath, strconv.Itoa(out)+".png")
		imaging.Save(img, outPath)
	}
}

func writeLyricsToImage(conf zazabul.Config, text string) image.Image {
	rootPath, _ := GetRootPath()

	fileHandle, err := os.Open(filepath.Join(rootPath, conf.Get("background_file")))
	if err != nil {
		panic(err)
	}
	pngData, _, err := image.Decode(fileHandle)
	if err != nil {
		panic(err)
	}
	b := pngData.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), pngData, b.Min, draw.Src)

	lyricsColor, _ := colorful.Hex(conf.Get("lyrics_color"))
	fg := image.NewUniform(lyricsColor)

	fontBytes, err := os.ReadFile(filepath.Join(rootPath, conf.Get("font_file")))
	if err != nil {
		panic(err)
	}
	fontParsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	c := freetype.NewContext()
	c.SetDPI(DPI)
	c.SetFont(fontParsed)
	c.SetFontSize(SIZE)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	texts := strings.Split(text, "\n")

	finalTexts := make([]string, 0)
	for _, txt := range texts {
		wrappedTxts := WordWrap(conf, txt, 1366-130)
		finalTexts = append(finalTexts, wrappedTxts...)
	}

	// Draw the text.
	pt := freetype.Pt(80, 50+int(c.PointToFixed(SIZE)>>6))
	for _, s := range finalTexts {
		_, err = c.DrawString(s, pt)
		if err != nil {
			panic(err)
		}
		pt.Y += c.PointToFixed(SIZE * SPACING)
	}

	return img
}
