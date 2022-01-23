package main

import (
  "os"
	color2 "github.com/gookit/color"
	"github.com/bankole7782/zazabul"
	"fmt"
  "time"
  "path/filepath"
  "image"
  "image/png"
  "image/color"
  "image/draw"
  "github.com/golang/freetype"
  "golang.org/x/image/font"
  "github.com/golang/freetype/truetype"
  "golang.org/x/image/math/fixed"
  "github.com/go-playground/colors"
  "bufio"
  "strconv"
  "strings"
  "os/exec"
)

const (
  DPI = 72.0
  SIZE = 80.0
  SPACING = 1.1
)

// 1366 - 130

func main() {
  rootPath, err := GetRootPath()
  if err != nil {
    panic(err)
    os.Exit(1)
  }

  if len(os.Args) < 2 {
		color2.Red.Println("Expecting a command. Run with help subcommand to view help.")
		os.Exit(1)
	}


	switch os.Args[1] {
	case "--help", "help", "h":
  		fmt.Println(`hananan is a terminal program that creates lyrics videos.
It uses a constant picture for the background.

It outputs frames which you would need to convert to video using ffmpeg.
The number of frames per seconds is 24. This is what this program uses.

Directory Commands:
  pwd     Print working directory. This is the directory where the files needed by any command
          in this cli program must reside.

Main Commands:
  init    Creates a config file describing your video. Edit to your own requirements.
          The file from init1 is expected for r1.

  run     Renders a project with the config created above. It expects a blender file and a
          launch file (created from 'init' above)
          All files must be placed in the working directory.

  pc      Print commands needed to convert the generated frames to video.

  			`)

  	case "pwd":
  		fmt.Println(rootPath)

    case "init":
      var	tmplOfMethod1 = `// output_name is the name of the project.
output_name:

// lyrics_file is the file that contains timestamps and lyrics chunks seperated by newlines.
// a sample can be found at https://sae.ng/static/bmtf.txt
//
lyrics_file:


// the font_file is the file of a ttf font that the text would be printed with.
// you could find a font on https://fonts.google.com
font_file:

// lyrics_color is the color of the rendered lyric. Example is #af1382
lyrics_color: #666666


// background_file is the background that would be used for this lyric video.
// the background_file must be a png
// the background_file must be of dimensions (1366 x 768)
background_file:

// total_length: The duration of the songs in this format (mm:ss)
total_length:

  	`
  		configFileName := "s" + time.Now().Format("20060102T150405") + ".zconf"
  		writePath := filepath.Join(rootPath, configFileName)

  		conf, err := zazabul.ParseConfig(tmplOfMethod1)
      if err != nil {
      	panic(err)
      }

      err = conf.Write(writePath)
      if err != nil {
        panic(err)
      }

      fmt.Printf("Edit the file at '%s' before launching.\n", writePath)


    case "run":
    	if len(os.Args) != 3 {
    		color2.Red.Println("The r1 command expects a file created by the init1 command")
    		os.Exit(1)
    	}

    	confPath := filepath.Join(rootPath, os.Args[2])

    	conf, err := zazabul.LoadConfigFile(confPath)
    	if err != nil {
    		panic(err)
    		os.Exit(1)
    	}

    	for _, item := range conf.Items {
    		if item.Value == "" {
    			color2.Red.Println("Every field in the launch file is compulsory.")
    			os.Exit(1)
    		}
    	}


      totalSeconds := timeFormatToSeconds(conf.Get("total_length"))
      lyricsObject := parseLyricsFile(filepath.Join(rootPath, conf.Get("lyrics_file")))
      renderPath := getRenderPath( conf.Get("output_name") )

      var lastSeconds int
      startedPrinting := false
      firstFrame := false

      for seconds := 0; seconds <= totalSeconds; seconds++ {

        if startedPrinting == false {
          _, ok := lyricsObject[seconds]
          if ! ok {
            fileHandle, err := os.Open(filepath.Join(rootPath, conf.Get("background_file")))
            if err != nil {
              panic(err)
            }
            img, _, err := image.Decode(fileHandle)
            if err != nil {
              panic(err)
            }
            writeImageToDisk(img, renderPath, seconds)
          } else {
            startedPrinting = true
            firstFrame = true
            lastSeconds = seconds
          }

        } else {

          img := writeToImage(conf, lyricsObject[lastSeconds])

          if firstFrame == true {
            writeImageToDisk(img, renderPath, lastSeconds )
            firstFrame = false
          }

          writeImageToDisk(img, renderPath, seconds)
          _, ok := lyricsObject[seconds]
          if ok {
            firstFrame = true
            lastSeconds = seconds
          }
        }

      }

      color2.Green.Println("Completed successfully. Output path: ", renderPath)

    case "pc":
      color2.Println("Switch to the folder created by the r1 command above.")
      color2.Green.Println("    ffmpeg -framerate 24 -i %d.png tmp_output.mp4")
      color2.Green.Println("    ffmpeg -i tmp_output.mp4 -i song.mp3 output.mp4")

    case "ff":
      out, err := exec.Command("$SNAP/bin/ffmpeg", "-h").CombinedOutput()
      if err != nil {
        panic(err)
      }
      fmt.Println(string(out))

  	default:
  		color2.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
  		os.Exit(1)
  	}

}


func writeImageToDisk(img image.Image, renderPath string, seconds int) {
  for i := 1; i <= 24; i++ {
    out := (24 * seconds) + i
    outPath := filepath.Join(renderPath, strconv.Itoa(out) + ".png")
    innerWriteImageToDisk(img, outPath)
  }
}

// Save that RGBA image to disk.
func innerWriteImageToDisk(img image.Image, outPath string) {
  outFile, err := os.Create(outPath)
  if err != nil {
    panic(err)
  }
  defer outFile.Close()
  b := bufio.NewWriter(outFile)
  err = png.Encode(b, img)
  if err != nil {
    panic(err)
  }
  err = b.Flush()
  if err != nil {
    panic(err)
  }
}


func writeToImage(conf zazabul.Config, text string) image.Image {
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

  hex, err := colors.ParseHEX(conf.Get("lyrics_color"))
  if err != nil {
    panic(err)
  }
  nCR := hex.ToRGBA()
  newColor := color.RGBA{uint8(nCR.R), uint8(nCR.G), uint8(nCR.B), 255}
  fg := image.NewUniform(newColor)


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
    wrappedTxts := wordWrap(conf, txt, 1366 - 130)
    finalTexts = append(finalTexts, wrappedTxts...)
  }

  if len(finalTexts) > 7 {
    color2.Red.Println("Shorten the following text for it to fit this video:")
    color2.Red.Println()
    for _, t := range strings.Split(text, "\n") {
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
