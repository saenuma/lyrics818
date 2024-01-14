package l8shared

import (
	"image"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const VersionFormat = "20060102T150405MST"

const (
	LAPTOP_WIDTH  = 1000
	LAPTOP_HEIGHT = 600
)

func MakeLaptopFrames(outName string, totalSeconds int, renderPath string, inputs map[string]string) {
	numberOfCPUS := runtime.NumCPU()
	lyricsObject := ParseLyricsFile(inputs["lyrics_file"], totalSeconds)

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
					bgColor, _ := colorful.Hex(inputs["background_color"])
					bg := image.NewUniform(bgColor)

					img := image.NewRGBA(image.Rect(0, 0, LAPTOP_WIDTH, LAPTOP_HEIGHT))
					draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

					outPath := filepath.Join(renderPath, strconv.Itoa(seconds)+".png")
					imaging.Save(img, outPath)
				} else {
					img := WriteLyricsToImage(inputs, lyricsObject[seconds])
					outPath := filepath.Join(renderPath, strconv.Itoa(seconds)+".png")
					imaging.Save(img, outPath)

				}

			}

		}(startSeconds, endSeconds, &wg)
	}
	wg.Wait()

	for seconds := (jobsPerThread * numberOfCPUS); seconds < totalSeconds; seconds++ {
		txt := lyricsObject[seconds]
		if txt == "" {
			bgColor, _ := colorful.Hex(inputs["background_color"])
			bg := image.NewUniform(bgColor)

			img := image.NewRGBA(image.Rect(0, 0, LAPTOP_WIDTH, LAPTOP_HEIGHT))
			draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

			outPath := filepath.Join(renderPath, strconv.Itoa(seconds)+".png")
			imaging.Save(img, outPath)
		} else {
			img := WriteLyricsToImage(inputs, lyricsObject[seconds])
			outPath := filepath.Join(renderPath, strconv.Itoa(seconds)+".png")
			imaging.Save(img, outPath)
		}
	}

}

func wordWrapLaptop(inputs map[string]string, text string, writeWidth int) []string {

	rgba := image.NewRGBA(image.Rect(0, 0, LAPTOP_WIDTH, LAPTOP_HEIGHT))

	fontBytes, err := os.ReadFile(inputs["font_file"])
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
			Size:    SIZE,
			DPI:     DPI,
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
			outStr := tmpStr[:len(tmpStr)-len(aStr)]
			tmpStr = oneStr
			outStrs = append(outStrs, outStr)
		}
	}
	outStrs = append(outStrs, tmpStr)

	return outStrs
}

func WriteLyricsToImage(inputs map[string]string, text string) image.Image {
	bgColor, _ := colorful.Hex(inputs["background_color"])
	bg := image.NewUniform(bgColor)

	img := image.NewRGBA(image.Rect(0, 0, LAPTOP_WIDTH, LAPTOP_HEIGHT))
	draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

	lyricsColor, _ := colorful.Hex(inputs["lyrics_color"])
	fg := image.NewUniform(lyricsColor)

	fontBytes, err := os.ReadFile(inputs["font_file"])
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
		wrappedTxts := wordWrapLaptop(inputs, txt, LAPTOP_WIDTH-160)
		finalTexts = append(finalTexts, wrappedTxts...)
	}

	// Draw the text.
	pt := freetype.Pt(60, 50+int(c.PointToFixed(SIZE)>>6))
	for _, s := range finalTexts {
		_, err = c.DrawString(s, pt)
		if err != nil {
			panic(err)
		}
		pt.Y += c.PointToFixed(SIZE * SPACING)
	}

	return img
}
