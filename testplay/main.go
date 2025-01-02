package main

import (
	"math"
	"os"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/internal"
)

func main() {
	if len(os.Args) < 2 {
		panic("expecting l8f file as only input")
	}

	songPath := os.Args[1]

	runtime.LockOSThread()

	internal.GetRootPath()
	ObjCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(1200, 800, "l8f format testplay", false)
	DrawNowPlayingUI(window, songPath, 0)

	// respond to the mouse
	// window.SetMouseButtonCallback(mouseBtnCallback)
	// window.SetCursorPosCallback(internal.CurPosCB)

	StartTime = time.Now()
	go playAudio(songPath)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		// update UI if song is playing
		if currentPlayer != nil && currentPlayer.IsPlaying() {
			seconds := time.Since(StartTime).Seconds()
			secondsInt := int(math.Floor(seconds))
			if secondsInt != CurrentPlaySeconds {
				DrawNowPlayingUI(window, songPath, secondsInt)
			}
		}

		time.Sleep(time.Second/time.Duration(internal.FPS) - time.Since(t))
	}
}
