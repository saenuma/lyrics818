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
	l8internal "github.com/saenuma/lyrics818/internal/lyrics818"
	"github.com/sqweek/dialog"
)

var colorObjCoords map[int]g143.Rect

func main() {
	rootPath, err := internal.GetRootPath()
	if err != nil {
		panic(err)
	}

	sampleLyricsPath := filepath.Join(rootPath, "bmtf.txt")
	os.WriteFile(sampleLyricsPath, internal.SampleLyricsFile, 0777)

	runtime.LockOSThread()

	l8internal.ObjCoords = make(map[int]g143.Rect)
	colorObjCoords = make(map[int]g143.Rect)
	l8internal.InputsStore = make(map[string]string)
	l8internal.InChannel = make(chan string)

	window := g143.NewWindow(1000, 800, "lyrics818: a more comfortable lyrics video generator", false)
	l8internal.AllDraws(window)

	go func() {
		for {
			method := <-l8internal.InChannel
			if method == "mp4" {
				ffPath := GetFFMPEGCommand()
				_, err := l8internal.MakeVideo(l8internal.InputsStore, ffPath)
				if err != nil {
					log.Println(err)
					return
				}

			} else if method == "l8f" {
				_, err := l8internal.MakeVideoL8F(l8internal.InputsStore)
				if err != nil {
					log.Println(err)
					return
				}
			}

			l8internal.ClearAfterRender = true
		}
	}()

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// respond to mouse movement
	window.SetCursorPosCallback(l8internal.CursorPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		if l8internal.ClearAfterRender {
			// clear the UI and redraw
			l8internal.InputsStore = make(map[string]string)
			l8internal.AllDraws(window)
			l8internal.DrawEndRenderView(window, l8internal.EmptyFrameNoInputs)
			time.Sleep(5 * time.Second)
			l8internal.AllDraws(window)

			// respond to the mouse
			window.SetMouseButtonCallback(mouseBtnCallback)
			// respond to mouse movement
			window.SetCursorPosCallback(l8internal.CursorPosCB)

			l8internal.ClearAfterRender = false
		}

		time.Sleep(time.Second/time.Duration(l8internal.FPS) - time.Since(t))
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

	for code, RS := range l8internal.ObjCoords {
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
	case l8internal.OpenWDBtn:
		l8internal.ExternalLaunch(rootPath)

	case l8internal.ViewLyricsBtn:
		sampleLyricsPath := filepath.Join(rootPath, "bmtf.txt")
		l8internal.ExternalLaunch(sampleLyricsPath)

	case l8internal.SelectLyricsBtn:
		filename, err := dialog.File().Filter("Lyrics File", "txt").Load()
		if filename == "" || err != nil {
			return
		}

		l8internal.InputsStore["lyrics_file"] = filename

		currentFrame := l8internal.RefreshInputsOnWindow(window, l8internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case l8internal.FontFileBtn:
		filename, err := dialog.File().Filter("Font file", "ttf").Load()
		if filename == "" || err != nil {
			return
		}
		l8internal.InputsStore["font_file"] = filename
		currentFrame := l8internal.RefreshInputsOnWindow(window, l8internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case l8internal.BgFileBtn:
		filename, err := dialog.File().Filter("PNG Image", "png").Load()
		if filename == "" || err != nil {
			return
		}

		l8internal.InputsStore["background_file"] = filename

		currentFrame := l8internal.RefreshInputsOnWindow(window, l8internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case l8internal.MusicFileBtn:
		filename, err := dialog.File().Filter("MP3 Audio", "mp3").Load()
		if filename == "" || err != nil {
			return
		}
		l8internal.InputsStore["music_file"] = filename

		currentFrame := l8internal.RefreshInputsOnWindow(window, l8internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case l8internal.LyricsColorBtn:
		drawPickColors(window)
		window.SetMouseButtonCallback(pickColorsMouseCallback)
		window.SetCursorPosCallback(nil)

	case l8internal.OurSite:
		l8internal.ExternalLaunch("https://sae.ng")

	case l8internal.RenderBtn:
		if len(l8internal.InputsStore) != 5 {
			return
		}

		currentFrame := l8internal.RefreshInputsOnWindow(window, l8internal.EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		l8internal.DrawRenderView(window, currentFrame)
		l8internal.InChannel <- "mp4"

	case l8internal.RenderL8fBtn:
		if len(l8internal.InputsStore) != 5 {
			return
		}

		currentFrame := l8internal.RefreshInputsOnWindow(window, l8internal.EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		l8internal.DrawRenderView(window, currentFrame)
		l8internal.InChannel <- "l8f"
	}

}
