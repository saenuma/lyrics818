package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/saenuma/lyrics818/l8f"
	"github.com/saenuma/lyrics818/l8shared"
	sDialog "github.com/sqweek/dialog"
)

func main() {
	ffplayCmd := GetFFPlayCommand()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	myApp := app.New()
	w := myApp.NewWindow("Test Videos made with Lyrics818")

	vidBox := container.NewVBox()

	vidFileLabel := widget.NewLabel("")
	getVidFileBtn := widget.NewButton("Select lyrics818 Video", func() {
		filename, err := sDialog.File().Filter("lyrics818 video", "l8f").Load()
		if err == nil {
			vidFileLabel.SetText(filename)
		}
	})

	widthSelect := widget.NewSelect([]string{"laptop", "mobile"}, nil)
	widthSelect.Selected = "laptop"

	startAtEntry := widget.NewEntry()
	startAtEntry.SetText("0:00")

	tmpWf64, tmpHf64 := 1366*0.8, 768*0.8
	laptopW, laptopH := int(tmpWf64), int(tmpHf64)

	videoImage := canvas.NewImageFromImage(nil)

	playTime := widget.NewLabel("0:00")
	totalLengthLabel := widget.NewLabel("0:00")

	rootPath, _ := l8shared.GetRootPath()
	tmpAudioPath := filepath.Join(rootPath, "tmp_audio.mp3")

	beginPlayAt := func(seek, inVideoPath, mode string) {
		startTime := time.Now()
		go func() {
			exec.CommandContext(ctx, ffplayCmd, tmpAudioPath, "-nodisp", "-ss", seek).Run()
		}()
		beginSeconds := l8shared.TimeFormatToSeconds(seek)

		if mode == "laptop" {
			currFrame, _ := l8f.ReadLaptopFrame(inVideoPath, 0)
			tmp := imaging.Fit(*currFrame, laptopW, laptopH, imaging.Lanczos)
			videoImage.Image = tmp
			videoImage.FillMode = canvas.ImageFillOriginal
			videoImage.Refresh()
		} else if mode == "mobile" {
			currFrame, _ := l8f.ReadMobileFrame(inVideoPath, 0)
			tmp := imaging.Fit(*currFrame, 400, 500, imaging.Lanczos)
			videoImage.Image = tmp
			videoImage.FillMode = canvas.ImageFillOriginal
			videoImage.Refresh()
		}

		// We can wait for the sound to finish playing using something like this
		for {
			seconds := time.Since(startTime).Seconds() + float64(beginSeconds)
			playTime.SetText(l8shared.SecondsToMinutes(int(seconds)))
			// currFrame, _ = l8f.ReadLaptopFrame(inVideoPath, int(seconds))
			// tmp := imaging.Fit(*currFrame, laptopW, laptopH, imaging.Lanczos)
			// videoImage.Image = tmp
			// videoImage.Refresh()

			if mode == "laptop" {
				currFrame, _ := l8f.ReadLaptopFrame(inVideoPath, int(seconds))
				tmp := imaging.Fit(*currFrame, laptopW, laptopH, imaging.Lanczos)
				videoImage.Image = tmp
				videoImage.Refresh()
			} else if mode == "mobile" {
				currFrame, _ := l8f.ReadMobileFrame(inVideoPath, int(seconds))
				tmp := imaging.Fit(*currFrame, 400, 500, imaging.Lanczos)
				videoImage.Image = tmp

				videoImage.Refresh()
			}

			time.Sleep(time.Second)
		}
	}

	playBtn := widget.NewButton("Play Lyrics818 Video", func() {
		if vidFileLabel.Text == "" {
			return
		}

		audioBytes, err := l8f.ReadAudio(vidFileLabel.Text)
		if err != nil {
			panic(err)
		}

		os.WriteFile(tmpAudioPath, audioBytes, 0777)

		videoLength, err := l8f.GetVideoLength(vidFileLabel.Text)
		if err != nil {
			panic(err)
		}

		totalLengthLabel.SetText(l8shared.SecondsToMinutes(videoLength))

		toolsBox := container.NewCenter(container.NewHBox(playTime, totalLengthLabel))

		vidBox.Add(container.NewPadded(videoImage))
		vidBox.Add(toolsBox)

		go beginPlayAt(startAtEntry.Text, vidFileLabel.Text, widthSelect.Selected)

	})

	formBox := container.NewVBox(
		container.NewHBox(widget.NewLabel("Lyrics818 Video File: "), getVidFileBtn),
		vidFileLabel,
		container.NewHBox(widget.NewLabel("Laptop or Mobile: "), widthSelect),
		container.NewHBox(widget.NewLabel("Start at: "), container.New(&l8shared.LongEntry{}, startAtEntry)),

		widget.NewSeparator(),
		playBtn,
	)

	windowBox := container.NewHBox(formBox, vidBox)
	w.SetOnClosed(func() {
		cancel()
	})
	w.SetContent(windowBox)
	w.Resize(fyne.NewSize(1200, 600))
	w.ShowAndRun()
}
