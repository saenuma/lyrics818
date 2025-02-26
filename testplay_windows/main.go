package main

import (
	"math"
	"os"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/internal"
	tpinternal "github.com/saenuma/lyrics818/internal/testplay"
)

func main() {
	if len(os.Args) < 2 {
		panic("expecting l8f file as only input")
	}

	songPath := os.Args[1]
	tpinternal.SongPath = songPath

	runtime.LockOSThread()

	internal.GetRootPath()
	tpinternal.ObjCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(1200, 800, "l8f format testplay", false)
	tpinternal.DrawNowPlayingUI(window, songPath, 0)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// window.SetCursorPosCallback(tpinternal.CurPosCB)

	tpinternal.StartTime = time.Now()
	go playAudio(songPath)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		// update UI if song is playing
		if currentPlayer != nil && currentPlayer.IsPlaying() {
			seconds := time.Since(tpinternal.StartTime).Seconds()
			secondsInt := int(math.Floor(seconds))
			if secondsInt != tpinternal.CurrentPlaySeconds {
				tpinternal.DrawNowPlayingUI(window, songPath, secondsInt)
			}
		}

		time.Sleep(time.Second/time.Duration(tpinternal.FPS) - time.Since(t))
	}
}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range tpinternal.ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	switch widgetCode {
	case tpinternal.SwitchViewBtn:
		if tpinternal.DeviceView == "laptop" {
			tpinternal.DeviceView = "mobile"
		} else if tpinternal.DeviceView == "mobile" {
			tpinternal.DeviceView = "laptop"
		}

		seconds := time.Since(tpinternal.StartTime).Seconds()
		secondsInt := int(math.Floor(seconds))
		tpinternal.DrawNowPlayingUI(window, tpinternal.SongPath, secondsInt)

	}
}
