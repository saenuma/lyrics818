package main

import (
	"log"

	"github.com/sqweek/dialog"
)

func PickImageFile() string {
	filename, err := dialog.File().Filter("PNG Image", "png").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func PickTxtFile() string {
	filename, err := dialog.File().Filter("Lyrics File", "txt").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func PickFontFile() string {
	filename, err := dialog.File().Filter("Font file", "ttf").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func PickMp3File() string {
	filename, err := dialog.File().Filter("MP3 Audio", "mp3").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}
