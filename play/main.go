package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/saenuma/lyrics818/l8f"
)

func main() {
	if len(os.Args) == 1 {
		panic("Expecting the filename of a video as only argument")
	}

	inVideoPath := os.Args[1]

	myApp := app.New()
	w := myApp.NewWindow("playing lyrics818 video - " + inVideoPath)

	blackImg := image.NewRGBA(image.Rect(0, 0, 1366, 768))
	draw.Draw(blackImg, blackImg.Bounds(), image.NewUniform(color.Black), blackImg.Bounds().Min, draw.Src)

	videoImage := canvas.NewImageFromImage(blackImg)
	videoImage.FillMode = canvas.ImageFillOriginal

	audioBytes, err := l8f.ReadAudio(inVideoPath)
	if err != nil {
		panic(err)
	}

	// Decode file
	decodedMp3, err := mp3.NewDecoder(bytes.NewReader(audioBytes))
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

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

	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(decodedMp3)

	playTime := widget.NewLabel("0:00")

	videoLength, err := l8f.GetVideoLength(inVideoPath)
	if err != nil {
		panic(err)
	}

	beginPlayAt := func(player oto.Player, seek string) {
		// Play starts playing the sound and returns without waiting for it (Play() is async).
		player.Play()
		startTime := time.Now()

		beginSeconds := TimeFormatToSeconds(seek)
		currFrame, _ := l8f.ReadLaptopFrame(inVideoPath, beginSeconds)
		videoImage.Image = *currFrame
		videoImage.Refresh()

		// We can wait for the sound to finish playing using something like this
		for player.IsPlaying() {
			seconds := time.Since(startTime).Seconds() + float64(beginSeconds)
			playTime.SetText(secondsToMinutes(int(seconds)))
			currFrame, _ = l8f.ReadLaptopFrame(inVideoPath, int(seconds))
			videoImage.Image = *currFrame
			videoImage.Refresh()

			time.Sleep(time.Second)
		}
	}

	playBtn := widget.NewButton("play/pause", func() {
		if player.IsPlaying() {
			player.Pause()
		} else {
			go beginPlayAt(player, playTime.Text)
		}
	})
	totalLengthLabel := widget.NewLabel(secondsToMinutes(videoLength))

	toolsBox := container.NewCenter(container.NewHBox(playBtn, playTime, totalLengthLabel))

	windowBox := container.NewVBox(container.NewCenter(videoImage), toolsBox)
	w.SetContent(windowBox)

	go beginPlayAt(player, "0:00")

	w.SetFixedSize(false)
	w.ShowAndRun()
}
