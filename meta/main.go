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
		fmt.Println(`lyrics818.meta is a terminal program that updates the metadata of a album folder
that contains songs encoded in l8f format.
The album art should be called 'art.png' and stored in the folder of an album.

Directory Commands:
  pwd     Print working directory. This is the directory where the files needed by any command
          in this cli program must reside.

Main Commands:
  init    init creates a config file that would be contain common information to be added to each file
          in a folder of l8f songs.
          Edit to your own requirements.

  run     Renders a project with the config created above. It expects a a config file generated from
          'init' command above.
          All files must be placed in the working directory.

  			`)

	case "pwd":
		fmt.Println(rootPath)

	case "init":
		var tmplOfMethod2 = `// group_name is the name of the group / company that made the song.
// in mp3 this could mean the artist
group_name:

// the name of the album this song belongs to.
album_name:

// the country which the group resides in . Example Nigeria
country:

// the year which this album was released by the group
year:

// class is a comma seperated list of tags given to the song
// it replaces the genre of mp3
// all the tags in this field must be oneline.
// Example is: heavy_instruments, female_vocals
class: 

  	`
		configFileName := "meta_init_" + time.Now().Format("20060102T150405") + ".zconf"
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
		if len(os.Args) != 4 {
			color2.Red.Println("The run command expects a file created by the init command and a folder")
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

		tagsMap := make(map[string]string)
		for _, item := range conf.Items {
			tagsMap[item.Name] = item.Value
		}

		fullFolderPath := filepath.Join(rootPath, os.Args[3])
		dirFIs, err := os.ReadDir(fullFolderPath)
		if err != nil {
			panic(err)
		}

		validFiles := make([]string, 0)
		for _, dirFI := range dirFIs {
			if strings.HasSuffix(dirFI.Name(), ".l8f") {
				validFiles = append(validFiles, dirFI.Name())
			}
		}

		tmpPath := filepath.Join(rootPath, ".tmp")
		tmpSongsPath := filepath.Join(rootPath, ".tmp_"+time.Now().Format("20060102T150405"))
		os.MkdirAll(tmpPath, 0777)
		os.MkdirAll(tmpSongsPath, 0777)
		for _, filename := range validFiles {
			fullSongPath := filepath.Join(fullFolderPath, filename)
			tmpSongPath := filepath.Join(tmpSongsPath, filename)
			err = l8f.UpdateMeta(fullSongPath, tagsMap, tmpPath, tmpSongPath)
			if err != nil {
				panic(err)
			}
		}

		os.RemoveAll(fullFolderPath)
		os.Rename(tmpSongsPath, fullFolderPath)
		os.RemoveAll(tmpSongsPath)

	default:
		color2.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
		os.Exit(1)
	}

}
