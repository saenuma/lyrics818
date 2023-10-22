package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/saenuma/lyrics818/l8f"
)

func main() {
	if runtime.GOOS == "linux" {
		hd, _ := os.UserHomeDir()
		dd := os.Getenv("SNAP_USER_DATA")
		if !strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) && dd != "" {
			outpath := filepath.Join(dd, ".asoundrc")
			os.WriteFile(outpath, SoundRC, 0666)
		}
	}

	os.Setenv("FYNE_THEME", "dark")
	if len(os.Args) == 1 {
		panic("Expecting the filename of a video as only argument")
	}

	inVideoPath := os.Args[1]
	mobilePlay := false
	if len(os.Args) == 3 && os.Args[2] == "mobile" {
		mobilePlay = true
	}

	myApp := app.New()
	w := myApp.NewWindow("playing lyrics818 video - " + inVideoPath)

	tmpWf64, tmpHf64 := 1366*0.8, 768*0.8
	laptopW, laptopH := int(tmpWf64), int(tmpHf64)
	var blackImg *image.RGBA
	if mobilePlay {
		blackImg = image.NewRGBA(image.Rect(0, 0, 480, 600))
	} else {
		blackImg = image.NewRGBA(image.Rect(0, 0, laptopW, laptopH))
	}
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
		var currFrame *image.Image
		if mobilePlay {
			currFrame, _ = l8f.ReadMobileFrame(inVideoPath, 0)
			tmp := imaging.Fit(*currFrame, 400, 500, imaging.Lanczos)
			videoImage.Image = tmp
		} else {
			currFrame, _ = l8f.ReadLaptopFrame(inVideoPath, 0)
			tmp := imaging.Fit(*currFrame, laptopW, laptopH, imaging.Lanczos)
			videoImage.Image = tmp
		}
		videoImage.Refresh()

		// We can wait for the sound to finish playing using something like this
		for player.IsPlaying() {
			seconds := time.Since(startTime).Seconds() + float64(beginSeconds)
			playTime.SetText(secondsToMinutes(int(seconds)))
			if mobilePlay {
				currFrame, _ = l8f.ReadMobileFrame(inVideoPath, int(seconds))
				tmp := imaging.Fit(*currFrame, 480, 600, imaging.Lanczos)
				videoImage.Image = tmp
			} else {
				currFrame, _ = l8f.ReadLaptopFrame(inVideoPath, int(seconds))
				tmp := imaging.Fit(*currFrame, laptopW, laptopH, imaging.Lanczos)
				videoImage.Image = tmp
			}
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

	// w.SetFixedSize(true)
	w.ShowAndRun()
}
