package main

import (
  "os"
	color2 "github.com/gookit/color"
	"github.com/saenuma/zazabul"
	"fmt"
  "time"
  "path/filepath"
  "os/exec"
  "runtime"
  "net/http"
  "strings"
  "io"
)

const (
  DPI = 72.0
  SIZE = 80.0
  SPACING = 1.1
)

const VersionFormat = "20060102T150405MST"

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


  if runtime.GOOS == "windows" {
    newVersionStr := ""
    resp, err := http.Get("https://sae.ng/static/wapps/lyrics818.txt")
    if err != nil {
      fmt.Println(err)
    }
    if err == nil {
      defer resp.Body.Close()
      body, err := io.ReadAll(resp.Body)
      if err == nil && resp.StatusCode == 200 {
        newVersionStr = string(body)
      }
    }

    newVersionStr = strings.TrimSpace(newVersionStr)
    currentVersionStr = strings.TrimSpace(currentVersionStr)

    hnv := false
    if newVersionStr != "" && newVersionStr != currentVersionStr {
      time1, err1 := time.Parse(VersionFormat, newVersionStr)
      time2, err2 := time.Parse(VersionFormat, currentVersionStr)

      if err1 == nil && err2 == nil && time2.Before(time1) {
        hnv = true
      }
    }

    if hnv == true {
      fmt.Println("lyrics818 has an update.")
      fmt.Println("please visit 'https://sae.ng/lyrics818' for update instructions." )
      fmt.Println()
    }

  }

	switch os.Args[1] {
	case "--help", "help", "h":
  		fmt.Println(`lyrics818 is a terminal program that creates lyrics videos.
It uses a constant picture for the background.

Directory Commands:
  pwd     Print working directory. This is the directory where the files needed by any command
          in this cli program must reside.

Main Commands:
  init    Creates a config file describing your video. Edit to your own requirements.
          The file from 'init' is expected for the 'run' command.

  run     Renders a project with the config created above. It expects a a config file generated from
          'init' command above.
          All files must be placed in the working directory.

  			`)

  	case "pwd":
  		fmt.Println(rootPath)

    case "init":
      var	tmplOfMethod1 = `// lyrics_file is the file that contains timestamps and lyrics chunks seperated by newlines.
// a sample can be found at https://sae.ng/static/bmtf.txt
lyrics_file:


// the font_file is the file of a ttf font that the text would be printed with.
// you could find a font on https://fonts.google.com
font_file:

// lyrics_color is the color of the rendered lyric. Example is #af1382
lyrics_color: #666666


// background_file is the background that would be used for this lyric video.
// the background_file must be a png or an mp4
// the background_file must be of dimensions (1366px x 768px)
// the framerate must be 60fps
// the length of the video must match the length of the song
// you can generate an mp4 from videos229
background_file:

// total_length: The duration of the songs in this format (mm:ss)
total_length:

// music_file is the song to add its audio to the video.
music_file:

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
    		color2.Red.Println("The run command expects a file created by the init command")
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


      outName := "s" + time.Now().Format("20060102T150405")
      totalSeconds := timeFormatToSeconds(conf.Get("total_length"))
      renderPath := filepath.Join(rootPath, outName)
      os.MkdirAll(renderPath, 0777)


      if filepath.Ext(conf.Get("background_file")) == ".png" {
        imageMethod(outName, totalSeconds, renderPath, conf)
      } else if filepath.Ext(conf.Get("background_file")) == ".mp4" {
        videoMethod(outName, totalSeconds, renderPath, conf)
      } else {
        color2.Red.Println("Unsupported backround_file format: must be .png or .mp4")
        os.Exit(1)
      }

      command := GetFFMPEGCommand()

      out, err := exec.Command(command, "-i", filepath.Join(renderPath, "tmp_" + outName + ".mp4"),
        "-i", filepath.Join(rootPath, conf.Get("music_file")), "-pix_fmt",  "yuv420p",
        filepath.Join(rootPath, outName + ".mp4") ).CombinedOutput()
      if err != nil {
        fmt.Println(string(out))
        panic(err)
      }

      // clearing temporary files
      os.RemoveAll(renderPath)

      color2.Green.Println("The video has been generated into: ", filepath.Join(rootPath, outName + ".mp4") )

  	default:
  		color2.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
  		os.Exit(1)
  	}

}
