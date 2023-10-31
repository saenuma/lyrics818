package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/saenuma/lyrics818/l8f"
	"github.com/saenuma/lyrics818/l8shared"
	sDialog "github.com/sqweek/dialog"
)

func main() {
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

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	samplingRate := 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	tmpWf64, tmpHf64 := 1366*0.8, 768*0.8
	laptopW, laptopH := int(tmpWf64), int(tmpHf64)

	blackImg := image.NewRGBA(image.Rect(0, 0, laptopW, laptopH))
	draw.Draw(blackImg, blackImg.Bounds(), image.NewUniform(color.Black), blackImg.Bounds().Min, draw.Src)

	videoImage := canvas.NewImageFromImage(blackImg)
	videoImage.FillMode = canvas.ImageFillOriginal

	playTime := widget.NewLabel("0:00")
	totalLengthLabel := widget.NewLabel("0:00")

	beginPlayAt := func(player oto.Player, seek, inVideoPath string) {
		// Play starts playing the sound and returns without waiting for it (Play() is async).
		player.Play()
		startTime := time.Now()

		beginSeconds := l8shared.TimeFormatToSeconds(seek)

		currFrame, _ := l8f.ReadLaptopFrame(inVideoPath, 0)
		tmp := imaging.Fit(*currFrame, laptopW, laptopH, imaging.Lanczos)
		videoImage.Image = tmp
		videoImage.Refresh()

		// We can wait for the sound to finish playing using something like this
		for player.IsPlaying() {
			seconds := time.Since(startTime).Seconds() + float64(beginSeconds)
			playTime.SetText(l8shared.SecondsToMinutes(int(seconds)))
			currFrame, _ = l8f.ReadLaptopFrame(inVideoPath, int(seconds))
			tmp := imaging.Fit(*currFrame, laptopW, laptopH, imaging.Lanczos)
			videoImage.Image = tmp
			videoImage.Refresh()

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

		videoLength, err := l8f.GetVideoLength(vidFileLabel.Text)
		if err != nil {
			panic(err)
		}

		totalLengthLabel.SetText(l8shared.SecondsToMinutes(videoLength))

		// Decode file
		decodedMp3, err := mp3.NewDecoder(bytes.NewReader(audioBytes))
		if err != nil {
			panic("mp3.NewDecoder failed: " + err.Error())
		}

		// Create a new 'player' that will handle our sound. Paused by default.
		player := otoCtx.NewPlayer(decodedMp3)

		toolsBox := container.NewCenter(container.NewHBox(playTime, totalLengthLabel))

		vidBox.Add(container.NewPadded(videoImage))
		vidBox.Add(toolsBox)

		go beginPlayAt(player, "0:00", vidFileLabel.Text)

	})

	formBox := container.NewVBox(
		container.NewHBox(widget.NewLabel("Lyrics818 Video File: "), getVidFileBtn),
		vidFileLabel,
		container.NewHBox(widget.NewLabel("Laptop or Mobile: "), widthSelect),

		widget.NewSeparator(),
		playBtn,
	)

	windowBox := container.NewHBox(formBox, vidBox)
	w.SetContent(windowBox)
	w.Resize(fyne.NewSize(1200, 600))
	w.ShowAndRun()
}
