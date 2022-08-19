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
	"github.com/otiai10/copy"
	"github.com/saenuma/zazabul"
	"golang.org/x/image/font"
)

func videoMethod(outName string, totalSeconds int, renderPath string, conf zazabul.Config) {
	rootPath, _ := GetRootPath()

	command := GetFFMPEGCommand()

	framesPath := filepath.Join(rootPath, "frames_"+outName)
	os.MkdirAll(framesPath, 0777)
	out, err := exec.Command(command, "-i", filepath.Join(rootPath, conf.Get("background_file")),
		"-r", "60/1", filepath.Join(framesPath, "%d.png")).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	color2.Green.Println("Finished getting frames from your video")

	lyricsObject := parseLyricsFile(filepath.Join(rootPath, conf.Get("lyrics_file")), totalSeconds)

	numberOfCPUS := runtime.NumCPU()

	jobsPerThread := int(math.Floor(float64(totalSeconds) * float64(60.0) / float64(numberOfCPUS)))

	// remainder := int(math.Mod(float64(totalSeconds), float64(numberOfCPUS)))
	var wg sync.WaitGroup

	for threadIndex := 0; threadIndex < numberOfCPUS; threadIndex++ {
		wg.Add(1)

		startFrame := threadIndex * jobsPerThread
		endFrame := (threadIndex + 1) * jobsPerThread

		go func(startFrame, endFrame int, wg *sync.WaitGroup) {
			defer wg.Done()
			for frameCount := startFrame; frameCount < endFrame; frameCount++ {
				seconds := frameCount / 60
				videoFramePath := filepath.Join(framesPath, strconv.Itoa(frameCount)+".png")

				txt := lyricsObject[seconds]
				if txt == "" {
					newPath := filepath.Join(renderPath, filepath.Base(videoFramePath))
					copy.Copy(videoFramePath, newPath)
				} else {
					img := writeLyricsToVideoFrame(conf, lyricsObject[seconds], videoFramePath)
					imaging.Save(img, filepath.Join(renderPath, strconv.Itoa(frameCount)+".png"))
				}

			}

		}(startFrame, endFrame, &wg)
	}
	wg.Wait()

	for frameCount := (jobsPerThread * numberOfCPUS); frameCount < totalSeconds*60; frameCount++ {
		seconds := frameCount / 60
		videoFramePath := filepath.Join(framesPath, strconv.Itoa(frameCount)+".png")

		txt := lyricsObject[seconds]
		if txt == "" {
			newPath := filepath.Join(renderPath, filepath.Base(videoFramePath))
			copy.Copy(videoFramePath, newPath)
		} else {
			img := writeLyricsToVideoFrame(conf, lyricsObject[seconds], videoFramePath)
			imaging.Save(img, filepath.Join(renderPath, strconv.Itoa(frameCount)+".png"))
		}
	}

	color2.Green.Println("Completed generating frames of your lyrics video")

	out, err = exec.Command(command, "-framerate", "60", "-i", filepath.Join(renderPath, "%d.png"),
		"-pix_fmt", "yuv420p",
		filepath.Join(renderPath, "tmp_"+outName+".mp4")).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	os.RemoveAll(framesPath)

}

func writeLyricsToVideoFrame(conf zazabul.Config, text, videoFramePath string) image.Image {
	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	pngData, err := imaging.Open(videoFramePath)
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

	texts := strings.Split(text, "\r\n")

	finalTexts := make([]string, 0)
	for _, txt := range texts {
		wrappedTxts := wordWrap(conf, txt, 1366-130)
		finalTexts = append(finalTexts, wrappedTxts...)
	}

	if len(finalTexts) > 7 {
		color2.Red.Println("Shorten the following text for it to fit this video:")
		color2.Red.Println()
		for _, t := range strings.Split(text, "\r\n") {
			color2.Red.Println("    ", t)
		}

		os.Exit(1)
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
