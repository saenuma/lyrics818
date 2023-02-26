package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	color2 "github.com/gookit/color"
	"github.com/saenuma/lyrics818/l8f"
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
  init    init creates a lyric video config file for a single image background and would produce
          an l8f file.
          Edit to your own requirements.

  run     Renders a project with the config created above. It expects a a config file generated from
          'init' command above.
          All files must be placed in the working directory.

  			`)

	case "pwd":
		fmt.Println(rootPath)

	case "init":
		var tmplOfMethod2 = `// lyrics_file is the file that contains timestamps and lyrics chunks seperated by newlines.
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
laptop_background_file:

// background_file is the background that would be used for this lyric video.
// the background_file must be a png
// the background_file must be of dimensions (800px x 1000px)
mobile_background_file:

// music_file is the song to add its audio to the video.
// lyrics818 expects a mp3 music file
// the music_file determines the duration of the video.
music_file:

  	`
		configFileName := "init_" + time.Now().Format("20060102T150405") + ".zconf"
		writePath := filepath.Join(rootPath, configFileName)

		conf, err := zazabul.ParseConfig(tmplOfMethod2)
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

		confFileName := os.Args[2]
		confPath := filepath.Join(rootPath, confFileName)

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

		fullMp3Path := filepath.Join(rootPath, conf.Get("music_file"))
		if !strings.HasSuffix(fullMp3Path, ".mp3") {
			color2.Red.Println("Expecting an mp3 file in 'music_file'")
			os.Exit(1)
		}

		totalSeconds, err := ReadSecondsFromMusicFile(fullMp3Path)
		if err != nil {
			panic(err)
		}

		laptopOutName := "lframes_" + time.Now().Format("20060102T150405")
		mobileOutName := "mframes_" + time.Now().Format("20060102T150405")
		lrenderPath := filepath.Join(rootPath, laptopOutName)
		os.MkdirAll(lrenderPath, 0777)
		mrenderPath := filepath.Join(rootPath, mobileOutName)
		os.MkdirAll(mrenderPath, 0777)

		makeLaptopFrames(laptopOutName, totalSeconds, lrenderPath, conf)
		makeMobileFrames(mobileOutName, totalSeconds, mrenderPath, conf)

		outName := "video_" + time.Now().Format("20060102T150405") + ".l8f"
		fullOutPath := filepath.Join(rootPath, outName)
		err = l8f.MakeL8F(lrenderPath, mrenderPath, fullMp3Path, map[string]string{"framerate": "24"},
			rootPath, fullOutPath)
		if err != nil {
			panic(err)
		}
		os.RemoveAll(lrenderPath)
		os.RemoveAll(mrenderPath)

		color2.Green.Println("The video has been generated into: ", fullOutPath)

	default:
		color2.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
		os.Exit(1)
	}

}
