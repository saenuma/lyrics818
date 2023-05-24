package main

import (
	"bytes"
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
		dialog.ShowCustom("About keys117", "Close", boxes, myWindow)
	})
	topBar := container.NewHBox(openWDBtn, viewSampleBtn)
	formBox := container.NewPadded()
	outputsBox := container.NewVBox()

	getLyricsForm := func() *widget.Form {

		txtFiles := getFilesOfType(rootPath, ".txt")
		mp3Files := getFilesOfType(rootPath, ".mp3")
		ttfFiles := getFilesOfType(rootPath, ".ttf")
		pngFiles := getFilesOfType(rootPath, ".png")

		lyricsInputForm := widget.NewForm()
		lyricsInputForm.Append("lyrics_file", widget.NewSelect(txtFiles, nil))
		lyricsInputForm.Append("font_file", widget.NewSelect(ttfFiles, nil))
		lyricsInputForm.Append("background_file", widget.NewSelect(pngFiles, nil))
		lyricsInputForm.Append("music_file", widget.NewSelect(mp3Files, nil))
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

			command := GetFFMPEGCommand()
			outFileName, err := l8shared.MakeVideo(inputs, command)
			if err != nil {
				log.Println(err)
				outputsBox.Add(widget.NewLabel("Error occured: " + err.Error()))
				return
			}
			openOutputButton := widget.NewButton("Open Video", func() {
				if runtime.GOOS == "windows" {
					exec.Command("cmd", "/C", "start", filepath.Join(rootPath, outFileName)).Run()
				} else if runtime.GOOS == "linux" {
					exec.Command("xdg-open", filepath.Join(rootPath, outFileName)).Run()
				}
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

func getFilesOfType(rootPath, ext string) []string {
	dirFIs, err := os.ReadDir(rootPath)
	if err != nil {
		panic(err)
	}
	files := make([]string, 0)
	for _, dirFI := range dirFIs {
		if !dirFI.IsDir() && !strings.HasPrefix(dirFI.Name(), ".") && strings.HasSuffix(dirFI.Name(), ext) {
			files = append(files, dirFI.Name())
		}

		if dirFI.IsDir() && !strings.HasPrefix(dirFI.Name(), ".") {
			innerDirFIs, _ := os.ReadDir(filepath.Join(rootPath, dirFI.Name()))
			innerFiles := make([]string, 0)

			for _, innerDirFI := range innerDirFIs {
				if !innerDirFI.IsDir() && !strings.HasPrefix(innerDirFI.Name(), ".") && strings.HasSuffix(innerDirFI.Name(), ext) {
					innerFiles = append(innerFiles, filepath.Join(dirFI.Name(), innerDirFI.Name()))
				}
			}

			if len(innerFiles) > 0 {
				files = append(files, innerFiles...)
			}
		}

	}

	return files
}
