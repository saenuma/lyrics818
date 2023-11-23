package main

import (
	"bytes"
	"image"
	"log"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/saenuma/lyrics818/l8shared"
	sDialog "github.com/sqweek/dialog"
)

func main() {
	rootPath, err := l8shared.GetRootPath()
	if err != nil {
		panic(err)
	}

	myApp := app.New()
	myApp.Settings().SetTheme(&l8shared.MyTheme{})

	myWindow := myApp.NewWindow("lyrics818: a more comfortable lyrics video generator")
	myWindow.SetOnClosed(func() {
	})

	openWDBtn := widget.NewButton("Open Outputs Directory", func() {
		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "start", rootPath).Run()
		} else if runtime.GOOS == "linux" {
			exec.Command("xdg-open", rootPath).Run()
		}
	})

	viewSampleBtn := widget.NewButton("View Sample Lyrics File", func() {
		sampleLyricsLabel := widget.NewLabel(string(l8shared.SampleLyricsFile))
		innerBox := container.New(&l8shared.FillSpace{}, container.NewMax(container.NewScroll(sampleLyricsLabel)))
		dialog.ShowCustom("Sample Lyrics File", "Close", innerBox, myWindow)
	})

	saeBtn := widget.NewButton("sae.ng", func() {
		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "start", "https://sae.ng").Run()
		} else if runtime.GOOS == "linux" {
			exec.Command("xdg-open", "https://sae.ng").Run()
		}
	})

	aboutBtn := widget.NewButton("About Us", func() {
		img, _, err := image.Decode(bytes.NewReader(l8shared.SaeLogoBytes))
		if err != nil {
			panic(err)
		}
		logoImage := canvas.NewImageFromImage(img)
		logoImage.FillMode = canvas.ImageFillOriginal

		boxes := container.NewVBox(
			container.NewCenter(logoImage),
			widget.NewLabelWithStyle("Brought to You with Love by", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			saeBtn,
		)
		dialog.ShowCustom("About lyrics818", "Close", boxes, myWindow)
	})
	topBar := container.NewCenter(container.NewHBox(openWDBtn, viewSampleBtn, aboutBtn))

	helpWidget := widget.NewRichTextFromMarkdown(`
## Help

1. Only .mp3 files are allowed for the **music file**	

1. Only .ttf files are allowed for the **font file**

1. Only .txt files are allowed for the **lyrics file**
	`)

	// formBox := container.NewPadded()
	outputsBox := container.NewVBox()

	lyricsFileLabel := widget.NewLabel("")
	getLyricsFileBtn := widget.NewButton("Get Lyrics File", func() {
		filename, err := sDialog.File().Filter("Lyrics file", "txt").Load()
		if err == nil {
			lyricsFileLabel.SetText(filename)
		}
	})

	mp3FileLabel := widget.NewLabel("")
	getMp3FileBtn := widget.NewButton("Get Mp3 File", func() {
		filename, err := sDialog.File().Filter("Mp3 Audio file", "mp3").Load()
		if err == nil {
			mp3FileLabel.SetText(filename)
		}
	})

	fontFileLabel := widget.NewLabel("")
	getFontFileBtn := widget.NewButton("Get Font ttf file", func() {
		filename, err := sDialog.File().Filter("Font ttf file", "ttf").Load()
		if err == nil {
			fontFileLabel.SetText(filename)
		}
	})

	colorEntry := widget.NewEntry()
	colorEntry.SetText("#666666")

	bgColorEntry := widget.NewEntry()
	bgColorEntry.SetText("#ffffff")

	makeButton := widget.NewButton("Make Lyrics Video", func() {

		inputs := map[string]string{
			"lyrics_file":      lyricsFileLabel.Text,
			"font_file":        fontFileLabel.Text,
			"music_file":       mp3FileLabel.Text,
			"lyrics_color":     colorEntry.Text,
			"background_color": bgColorEntry.Text,
		}

		for _, v := range inputs {
			if v == "" {
				return
			}
		}

		outLabel := widget.NewLabel("Beginning")
		outputsBox.Add(outLabel)

		_, err := l8shared.MakeVideo2(inputs)
		if err != nil {
			log.Println(err)
			outputsBox.Add(widget.NewLabel("Error occured: " + err.Error()))
			return
		}
		outLabel.SetText("Done")
		outputsBox.Refresh()
	})
	makeButton.Importance = widget.HighImportance

	formBox := container.NewVBox(
		container.NewHBox(widget.NewLabel("Music File: "), getMp3FileBtn, mp3FileLabel),
		container.NewHBox(widget.NewLabel("Lyrics File: "), getLyricsFileBtn, lyricsFileLabel),
		container.NewHBox(widget.NewLabel("Background Color: "), container.New(&l8shared.LongEntry{}, bgColorEntry)),
		container.NewHBox(widget.NewLabel("Font File: "), getFontFileBtn, fontFileLabel),
		container.NewHBox(widget.NewLabel("Lyrics Color: "), container.New(&l8shared.LongEntry{}, colorEntry)),

		widget.NewSeparator(),
		container.New(&l8shared.Halfes{}, makeButton),
	)

	guitarImg, _, err := image.Decode(bytes.NewReader(l8shared.GuitarJPG))
	if err != nil {
		panic(err)
	}
	guitarFyneImage := canvas.NewImageFromImage(guitarImg)
	guitarFyneImage.FillMode = canvas.ImageFillOriginal
	guitarBox := container.NewCenter(guitarFyneImage)

	windowBox := container.NewHBox(
		guitarBox,
		container.NewVBox(
			container.NewCenter(topBar),
			helpWidget, formBox, outputsBox,
		),
	)

	myWindow.SetContent(windowBox)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.ShowAndRun()
}
