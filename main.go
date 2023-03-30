package main

import (
	"bytes"
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// os.Setenv("FYNE_SCALE", "0.9")
	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})

	myWindow := myApp.NewWindow("lyrics818: a more comfortable lyrics video generator")
	myWindow.SetOnClosed(func() {
	})

	openWDBtn := widget.NewButton("Open Working Directory", func() {
		exec.Command("cmd", "/C", "start", rootPath).Run()
	})

	viewSampleBtn := widget.NewButton("View Sample Lyrics File", func() {
		sampleLyricsLabel := widget.NewLabel(string(sampleLyricsFile))
		innerBox := container.New(&fillSpace{}, container.NewMax(container.NewScroll(sampleLyricsLabel)))
		dialog.ShowCustom("Sample Lyrics File", "Close", innerBox, myWindow)
	})

	saeBtn := widget.NewButton("sae.ng", func() {
		exec.Command("cmd", "/C", "start", "https://sae.ng").Run()
	})

	aboutBtn := widget.NewButton("About Us", func() {
		img, _, err := image.Decode(bytes.NewReader(SaeLogoBytes))
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
		dialog.ShowCustom("About keys117", "Close", boxes, myWindow)
	})
	topBar := container.NewHBox(openWDBtn, viewSampleBtn)
	formBox := container.NewPadded()
	outputsBox := container.NewVBox()

	getLyricsForm := func() *widget.Form {
		dirFIs, err := os.ReadDir(rootPath)
		if err != nil {
			panic(err)
		}
		files := make([]string, 0)
		for _, dirFI := range dirFIs {
			if !dirFI.IsDir() && !strings.HasPrefix(dirFI.Name(), ".") {
				files = append(files, dirFI.Name())
			}
		}

		lyricsInputForm := widget.NewForm()
		lyricsInputForm.Append("lyrics_file", widget.NewSelect(files, nil))
		lyricsInputForm.Append("font_file", widget.NewSelect(files, nil))
		lyricsInputForm.Append("background_file", widget.NewSelect(files, nil))
		lyricsInputForm.Append("music_file", widget.NewSelect(files, nil))
		colorEntry := widget.NewEntry()
		colorEntry.SetText("#666666")
		lyricsInputForm.Append("lyrics_color", colorEntry)
		lyricsInputForm.SubmitText = "Make Lyrics Video"
		lyricsInputForm.CancelText = "Close"
		lyricsInputForm.OnCancel = func() {
			os.Exit(0)
		}
		lyricsInputForm.OnSubmit = func() {
			outputsBox.Add(widget.NewLabel("Beginning"))
			inputs := getFormInputs(lyricsInputForm.Items)
			outFileName, err := makeLyrics(inputs)
			if err != nil {
				log.Println(err)
				outputsBox.Add(widget.NewLabel("Error occured: " + err.Error()))
				return
			}
			openOutputButton := widget.NewButton("Open Video", func() {
				exec.Command("cmd", "/C", "start", filepath.Join(rootPath, outFileName)).Run()
			})
			outputsBox.Add(openOutputButton)
			outputsBox.Refresh()
		}

		return lyricsInputForm
	}

	refreshBtn := widget.NewButton("Refresh Files List", func() {
		lyricsForm := getLyricsForm()
		formBox.RemoveAll()
		formBox.Add(lyricsForm)
		formBox.Refresh()
	})

	topBar.Add(refreshBtn)
	topBar.Add(aboutBtn)
	helpWidget := widget.NewRichTextFromMarkdown(`
## Help
1. All files must be placed in the working directory of this program.

1. Only .mp3 files are allowed for the **input music file**	

1. Only .png files are allowed for the **background**

1. The background_file must be of dimensions (1366px x 768px)
	`)
	windowBox := container.NewVBox(
		topBar,
		widget.NewSeparator(),
		helpWidget,
		formBox, outputsBox,
	)

	lyricsForm := getLyricsForm()
	formBox.Add(lyricsForm)
	formBox.Refresh()

	myWindow.SetContent(windowBox)
	myWindow.Resize(fyne.NewSize(800, 600))
	// myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}

func getFormInputs(content []*widget.FormItem) map[string]string {
	inputs := make(map[string]string)
	for _, formItem := range content {
		entryWidget, ok := formItem.Widget.(*widget.Entry)
		if ok {
			inputs[formItem.Text] = entryWidget.Text
			continue
		}

		selectWidget, ok := formItem.Widget.(*widget.Select)
		if ok {
			inputs[formItem.Text] = selectWidget.Selected
		}
	}

	return inputs
}
