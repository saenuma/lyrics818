package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/internal"
)

func main() {
	rootPath, err := internal.GetRootPath()
	if err != nil {
		panic(err)
	}

	sampleLyricsPath := filepath.Join(rootPath, "bmtf.txt")
	os.WriteFile(sampleLyricsPath, internal.SampleLyricsFile, 0777)

	runtime.LockOSThread()

	internal.ObjCoords = make(map[int]g143.Rect)
	internal.InputsStore = make(map[string]string)
	internal.InChannel = make(chan bool)
	internal.InColorChannel = make(chan bool)

	wWidth, wHeight := 1000, 800
	window := g143.NewWindow(wWidth, wHeight, "lyrics818: a more comfortable lyrics video generator", false)
	internal.AllDraws(window)

	go func() {
		for {
			<-internal.InChannel

			ffPath := GetFFMPEGCommand()
			_, err := internal.MakeVideo(internal.InputsStore, ffPath)
			if err != nil {
				log.Println(err)
				return
			}
			internal.DoneWithRender = true
		}
	}()

	go func() {
		for {
			<-internal.InColorChannel
			internal.PickedColor = pickColor()
			internal.ClearAfterColorPicker = true
		}
	}()

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// respond to mouse movement
	window.SetCursorPosCallback(internal.CursorPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		if internal.DoneWithRender {
			// clear the UI and redraw
			internal.AllDraws(window)
			msg := "Done! Open working directory"
			currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
			internal.DrawDialogWithMessage(window, currentFrame, msg)

			time.Sleep(5 * time.Second)

			// send the frame to glfw window
			windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
			window.SwapBuffers()

			// respond to the mouse
			window.SetMouseButtonCallback(mouseBtnCallback)
			// respond to mouse movement
			window.SetCursorPosCallback(internal.CursorPosCB)

			internal.DoneWithRender = false
		}

		if internal.ClearAfterColorPicker {
			internal.InputsStore["lyrics_color"] = internal.PickedColor

			currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
			// send the frame to glfw window
			windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
			window.SwapBuffers()

			// respond to the mouse
			window.SetMouseButtonCallback(mouseBtnCallback)
			// respond to mouse movement
			window.SetCursorPosCallback(internal.CursorPosCB)

			internal.ClearAfterColorPicker = false
		}

		time.Sleep(time.Second/time.Duration(internal.FPS) - time.Since(t))
	}

}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range internal.ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	rootPath, _ := internal.GetRootPath()

	switch widgetCode {
	case internal.OpenWDBtn:
		internal.ExternalLaunch(rootPath)

	case internal.ViewLyricsBtn:
		sampleLyricsPath := filepath.Join(rootPath, "bmtf.txt")
		internal.ExternalLaunch(sampleLyricsPath)

	case internal.SelectLyricsBtn:
		filename := pickFileUbuntu("txt")
		if filename == "" {
			return
		}
		internal.InputsStore["lyrics_file"] = filename

		currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.FontFileBtn:
		filename := pickFileUbuntu("ttf")
		if filename == "" {
			return
		}
		internal.InputsStore["font_file"] = filename
		currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.BgFileBtn:
		filename := pickFileUbuntu("png")
		if filename == "" {
			return
		}
		internal.InputsStore["background_file"] = filename

		currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.MusicFileBtn:
		filename := pickFileUbuntu("mp3")
		if filename == "" {
			return
		}
		internal.InputsStore["music_file"] = filename

		currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.LyricsColorBtn:
		currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		internal.DrawDialogWithMessage(window, currentFrame, "Awaiting your color choice!")
		internal.InColorChannel <- true

	case internal.OurSite:
		internal.ExternalLaunch("https://sae.ng")

	case internal.RenderBtn:
		if len(internal.InputsStore) != 5 {
			return
		}

		currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		internal.DrawDialogWithMessage(window, currentFrame, "Rendering! Please Wait")
		internal.InChannel <- true

	}
}
