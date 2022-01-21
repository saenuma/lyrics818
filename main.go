package main

import (
  "os"
	"github.com/gookit/color"
	"github.com/bankole7782/zazabul"
	"fmt"
  "time"
  "path/filepath"
)


func main() {
  rootPath, err := GetRootPath()
  if err != nil {
    panic(err)
    os.Exit(1)
  }

  if len(os.Args) < 2 {
		color.Red.Println("Expecting a command. Run with help subcommand to view help.")
		os.Exit(1)
	}


	switch os.Args[1] {
	case "--help", "help", "h":
  		fmt.Println(`hananan is a terminal program that creates lyrics videos.
It outputs frames which you would need to convert to video using ffmpeg.

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


// background_file is the background that would be used for this lyric video.
background_file:

// total_length: The duration of the songs in this format (mm:ss)
total_length:

  	`
  		configFileName := "s" + time.Now().Format("20060102") + ".zconf"
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


    case "r":
    	if len(os.Args) != 3 {
    		color.Red.Println("The r1 command expects a file created by the init1 command")
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
    			color.Red.Println("Every field in the launch file is compulsory.")
    			os.Exit(1)
    		}
    	}



  	default:
  		color.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
  		os.Exit(1)
  	}

}
