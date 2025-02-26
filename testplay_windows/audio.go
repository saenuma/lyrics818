package main

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"github.com/saenuma/lyrics818/internal"
	"github.com/saenuma/lyrics818/l8f"
)

var GlobalOtoCtx *oto.Context
var currentPlayer *oto.Player

func playAudio(l8fPath string) {
	rootPath, _ := internal.GetRootPath()

	audioBytes, err := l8f.ReadAudio(l8fPath)
	if err != nil {
		panic(err)
	}

	tmpAudioPath := filepath.Join(rootPath, ".tmp_audio.mp3")
	os.WriteFile(tmpAudioPath, audioBytes, 0777)

	// Read the mp3 file into memory
	fileBytes, err := os.ReadFile(tmpAudioPath)
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	if GlobalOtoCtx == nil {
		op := &oto.NewContextOptions{}

		// Usually 44100 or 48000. Other values might cause distortions in Oto
		op.SampleRate = 44100

		// Number of channels (aka locations) to play sounds from. Either 1 or 2.
		// 1 is mono sound, and 2 is stereo (most speakers are stereo).
		op.ChannelCount = 2

		// Format of the source. go-mp3's format is signed 16bit integers.
		op.Format = oto.FormatSignedInt16LE

		// Remember that you should **not** create more than one context
		otoCtx, readyChan, err := oto.NewContext(op)
		if err != nil {
			panic("oto.NewContext failed: " + err.Error())
		}
		GlobalOtoCtx = otoCtx
		// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
		<-readyChan
	}

	// Create a new 'player' that will handle our sound. Paused by default.
	player := GlobalOtoCtx.NewPlayer(decodedMp3)
	currentPlayer = player
	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Second)
	}

}

func continueAudio() {
	currentPlayer.Play()
	// We can wait for the sound to finish playing using something like this
	for currentPlayer.IsPlaying() {
		time.Sleep(time.Second)
	}
}
