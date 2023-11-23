package main

import (
	"bytes"
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/saenuma/lyrics818/l8shared"
)

func main() {
	// os.Setenv("FYNE_SCALE", "0.9")
	rootPath, err := l8shared.GetRootPath()
	if err != nil {
		panic(err)
	}

	myApp := app.New()
	myApp.Settings().SetTheme(&l8shared.MyTheme{})

	myWindow := myApp.NewWindow("lyrics818: a more comfortable lyrics video generator")
	myWindow.SetOnClosed(func() {
	})

	openWDBtn := widget.NewButton("Open Working Directory", func() {
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
	topBar := container.NewHBox(openWDBtn, viewSampleBtn)
	formBox := container.NewPadded()
	outputsBox := container.NewVBox()

	getLyricsForm := func() *widget.Form {

		txtFiles := l8shared.GetFilesOfType(rootPath, ".txt")
		mp3Files := l8shared.GetFilesOfType(rootPath, ".mp3")
		ttfFiles := l8shared.GetFilesOfType(rootPath, ".ttf")

		bgColorEntry := widget.NewEntry()
		bgColorEntry.SetText("#ffffff")

		lyricsInputForm := widget.NewForm()
		lyricsInputForm.Append("music_file", widget.NewSelect(mp3Files, nil))
		lyricsInputForm.Append("lyrics_file", widget.NewSelect(txtFiles, nil))
		lyricsInputForm.Append("background_color", bgColorEntry)
		lyricsInputForm.Append("font_file", widget.NewSelect(ttfFiles, nil))

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

			// complete the paths
			for k, v := range inputs {
				if k == "lyrics_color" || k == "background_color" {
					continue
				} else {
					inputs[k] = filepath.Join(rootPath, v)
				}
			}

			_, err := l8shared.MakeVideo2(inputs)
			if err != nil {
				log.Println(err)
				outputsBox.Add(widget.NewLabel("Error occured: " + err.Error()))
				return
			}

			outputsBox.Add(widget.NewLabel("Done. Check Working Directory"))
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

1. Only .mp3 files are allowed for the **music file**	

1. Only .ttf files are allowed for the **font file**

1. Only .txt files are allowed for the **lyrics file**

`)
	rightBox := container.NewVBox(
		topBar,
		widget.NewSeparator(),
		helpWidget,
		formBox, outputsBox,
	)

	lyricsForm := getLyricsForm()
	formBox.Add(lyricsForm)
	formBox.Refresh()

	guitarImg, _, err := image.Decode(bytes.NewReader(l8shared.GuitarJPG))
	if err != nil {
		panic(err)
	}
	guitarFyneImage := canvas.NewImageFromImage(guitarImg)
	guitarFyneImage.FillMode = canvas.ImageFillOriginal
	guitarBox := container.NewCenter(guitarFyneImage)

	windowBox := container.NewHBox(guitarBox, rightBox)
	myWindow.SetContent(windowBox)
	myWindow.Resize(fyne.NewSize(1000, 600))
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
