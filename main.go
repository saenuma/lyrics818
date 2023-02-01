package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	color2 "github.com/gookit/color"
	"github.com/saenuma/zazabul"
)

const VersionFormat = "20060102T150405MST"

func main() {

	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		color2.Red.Println("Expecting a command. Run with help subcommand to view help.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--help", "help", "h":
		fmt.Println(`lyrics818 is a terminal program that creates lyrics videos.
It uses a constant picture for the background.

Directory Commands:
  pwd     Print working directory. This is the directory where the files needed by any command
          in this cli program must reside.

Main Commands:
  init1   Method 1 creates a lyric video config file for a single image background.
          Edit to your own requirements.

  init2   Method 2 creates a lyric video config file for multiple image backgrounds.
          Edit to your own requirements.

  run     Renders a project with the config created above. It expects a a config file generated from
          'init' command above.
          All files must be placed in the working directory.

  			`)

	case "pwd":
		fmt.Println(rootPath)

	case "init1":
		var tmplOfMethod1 = `// lyrics_file is the file that contains timestamps and lyrics chunks seperated by newlines.
// a sample can be found at https://sae.ng/static/bmtf.txt
lyrics_file:


// the font_file is the file of a ttf font that the text would be printed with.
// you could find a font on https://fonts.google.com
font_file:

// lyrics_color is the color of the rendered lyric. Example is #af1382
lyrics_color: #666666


// background_file is the background that would be used for this lyric video.
// the background_file must be a png
// the background_file must be of dimensions (1366px x 768px)
background_file:

// music_file is the song to add its audio to the video.
// lyrics818 expects a mp3 music file
// the music_file determines the duration of the video.
music_file:

  	`
		configFileName := "m1_" + time.Now().Format("20060102T150405") + ".zconf"
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

	case "init2":
		var tmplOfMethod1 = `// lyrics_file is the file that contains timestamps and lyrics chunks seperated by newlines.
// a sample can be found at https://sae.ng/static/bmtf.txt
lyrics_file:


// the font_file is the file of a ttf font that the text would be printed with.
// you could find a font on https://fonts.google.com
font_file:

// lyrics_color is the color of the rendered lyric. Example is #af1382
lyrics_color: #666666

// The directory containing the pictures for a slideshow. It must be stored in the working directory
// of lyrics818.
// All pictures here must be of width 1366px and height 768px
// the background_files must be png
pictures_dir:

// music_file is the song to add its audio to the video.
// lyrics818 expects a mp3 music file
// the music_file determines the duration of the video.
music_file:

  	`
		configFileName := "m2_" + time.Now().Format("20060102T150405") + ".zconf"
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
			color2.Red.Println("The run command expects a file created by the init command")
			os.Exit(1)
		}

		confPath := filepath.Join(rootPath, os.Args[2])

		conf, err := zazabul.LoadConfigFile(confPath)
		if err != nil {
			panic(err)
		}

		for _, item := range conf.Items {
			if item.Value == "" {
				color2.Red.Println("Every field in the launch file is compulsory.")
				os.Exit(1)
			}
		}

		outName := "s" + time.Now().Format("20060102T150405")

		fullMp3Path := filepath.Join(rootPath, conf.Get("music_file"))
		if !strings.HasSuffix(fullMp3Path, ".mp3") {
			color2.Red.Println("Expecting an mp3 file in 'music_file'")
			os.Exit(1)
		}

		totalSeconds, err := ReadSecondsFromMusicFile(fullMp3Path)
		if err != nil {
			panic(err)
		}

		renderPath := filepath.Join(rootPath, outName)
		os.MkdirAll(renderPath, 0777)

		command := GetFFMPEGCommand()

		if strings.HasPrefix(confPath, "m1_") {
			// run method 1
			Method1(outName, totalSeconds, renderPath, conf)
		} else if strings.HasPrefix(confPath, "m2_") {
			// run method 2
			Method2(outName, totalSeconds, renderPath, conf)
		} else {
			color2.Red.Println("Invalid lyrics818 config file")
			os.Exit(1)
		}

		out, err := exec.Command(command, "-i", filepath.Join(renderPath, "tmp_"+outName+".mp4"),
			"-i", filepath.Join(rootPath, conf.Get("music_file")), "-pix_fmt", "yuv420p",
			filepath.Join(rootPath, outName+".mp4")).CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			panic(err)
		}

		// clearing temporary files
		os.RemoveAll(renderPath)

		color2.Green.Println("The video has been generated into: ", filepath.Join(rootPath, outName+".mp4"))

	default:
		color2.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
		os.Exit(1)
	}

}
