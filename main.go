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
  "github.com/go-playground/colors"
  "bufio"
  "strconv"
  "strings"
)

const (
  DPI = 72.0
  SIZE = 45.0
  SPACING = 1.1
)


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
It outputs frames which you would need to convert to video using ffmpeg.
The number of frames per seconds is 24. This is what this program uses.

Directory Commands:
  pwd     Print working directory. This is the directory where the files needed by any command
          in this cli program must reside.

Method 1: This uses a constant picture for the background.

  init1   Creates a config file describing your video. Edit to your own requirements.
          The file from init1 is expected for r1.

  r1      Renders a project with the config created above. It expects a blender file and a
          launch file (created from 'init' above)
          All files must be placed in the working directory.
  			`)

  	case "pwd":
  		fmt.Println(rootPath)

    case "init1":
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


    case "r1":
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
      lastFrameCount := 1

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

  // Draw the text.
  pt := freetype.Pt(80, 50+int(c.PointToFixed(SIZE)>>6))
  for _, s := range texts {
    _, err = c.DrawString(s, pt)
    if err != nil {
      panic(err)
    }
    pt.Y += c.PointToFixed(SIZE * SPACING)
  }

  return img
}
