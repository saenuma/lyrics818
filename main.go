package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	sampleLyricsPath := filepath.Join(rootPath, "bmtf.txt")
	os.WriteFile(sampleLyricsPath, SampleLyricsFile, 0777)

	runtime.LockOSThread()

	ObjCoords = make(map[int]g143.Rect)
	InputsStore = make(map[string]string)
	InChannel = make(chan string)
	InColorChannel = make(chan bool)

	wWidth, wHeight := 700, 650
	window := g143.NewWindow(wWidth, wHeight, "lyrics818: a more comfortable lyrics video generator", false)
	AllDraws(window)

	go func() {
		for {
			method := <-InChannel
			if method == "mp4" {
				ffPath := GetFFMPEGCommand()
				_, err := MakeVideo(InputsStore, ffPath)
				if err != nil {
					log.Println(err)
					return
				}

			} else if method == "l8f" {
				_, err := MakeVideoL8F(InputsStore)
				if err != nil {
					log.Println(err)
					return
				}
			}

			DoneWithRender = true
		}
	}()

	go func() {
		for {
			<-InColorChannel
			PickedColor = pickColor()
			ClearAfterColorPicker = true
		}
	}()

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// respond to mouse movement
	window.SetCursorPosCallback(CursorPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		if DoneWithRender {
			// clear the UI and redraw
			AllDraws(window)
			msg := "Done! Open working directory"
			currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
			DrawDialogWithMessage(window, currentFrame, msg)

			time.Sleep(5 * time.Second)

			// send the frame to glfw window
			windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
			window.SwapBuffers()

			// respond to the mouse
			window.SetMouseButtonCallback(mouseBtnCallback)
			// respond to mouse movement
			window.SetCursorPosCallback(CursorPosCB)

			DoneWithRender = false
		}

		if ClearAfterColorPicker {
			InputsStore["lyrics_color"] = PickedColor

			currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
			// send the frame to glfw window
			windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
			window.SwapBuffers()

			// respond to the mouse
			window.SetMouseButtonCallback(mouseBtnCallback)
			// respond to mouse movement
			window.SetCursorPosCallback(CursorPosCB)

			ClearAfterColorPicker = false
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
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

	for code, RS := range ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	rootPath, _ := GetRootPath()

	switch widgetCode {
	case OpenWDBtn:
		ExternalLaunch(rootPath)

	case ViewLyricsBtn:
		sampleLyricsPath := filepath.Join(rootPath, "bmtf.txt")
		ExternalLaunch(sampleLyricsPath)

	case SelectLyricsBtn:
		filename := PickTxtFile()
		if filename == "" {
			return
		}

		InputsStore["lyrics_file"] = filename

		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case FontFileBtn:
		filename := PickFontFile()
		if filename == "" {
			return
		}
		InputsStore["font_file"] = filename
		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case BgFileBtn:
		filename := PickImageFile()
		if filename == "" {
			return
		}

		InputsStore["background_file"] = filename

		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case MusicFileBtn:
		filename := PickMp3File()
		if filename == "" {
			return
		}
		InputsStore["music_file"] = filename

		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case LyricsColorBtn:
		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		DrawDialogWithMessage(window, currentFrame, "Awaiting your color choice!")
		InColorChannel <- true

	case OurSite:
		ExternalLaunch("https://sae.ng")

	case RenderBtn:
		if len(InputsStore) != 5 {
			return
		}

		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		DrawDialogWithMessage(window, currentFrame, "Rendering! Please Wait")
		InChannel <- "mp4"

	case RenderL8fBtn:
		if len(InputsStore) != 5 {
			return
		}

		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		DrawDialogWithMessage(window, currentFrame, "Rendering! Please Wait")
		InChannel <- "l8f"

	}

}
