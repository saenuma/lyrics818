package main

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
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

	internal.ObjCoords = make(map[int]g143.RectSpecs)
	internal.InputsStore = make(map[string]string)
	internal.InChannel = make(chan string)

	window := g143.NewWindow(1000, 800, "lyrics818: a more comfortable lyrics video generator", false)
	internal.AllDraws(window)

	go func() {
		for {
			method := <-internal.InChannel
			if method == "mp4" {
				ffPath := GetFFMPEGCommand()
				_, err := internal.MakeVideo(internal.InputsStore, ffPath)
				if err != nil {
					log.Println(err)
					return
				}

			} else if method == "l8f" {
				_, err := internal.MakeVideoL8F(internal.InputsStore)
				if err != nil {
					log.Println(err)
					return
				}
			}

			internal.ClearAfterRender = true
		}
	}()

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// respond to mouse movement
	window.SetCursorPosCallback(cursorPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		if internal.ClearAfterRender {
			// clear the UI and redraw
			internal.InputsStore = make(map[string]string)
			internal.AllDraws(window)
			internal.DrawEndRenderView(window, internal.EmptyFrameNoInputs)
			time.Sleep(5 * time.Second)
			internal.AllDraws(window)

			// respond to the mouse
			window.SetMouseButtonCallback(mouseBtnCallback)
			// respond to mouse movement
			window.SetCursorPosCallback(cursorPosCB)

			internal.ClearAfterRender = false
		}

		time.Sleep(time.Second/time.Duration(internal.FPS) - time.Since(t))
	}

}

func refreshInputsOnWindow(window *glfw.Window, frame image.Image) image.Image {
	wWidth, _ := window.GetSize()

	ggCtx := gg.NewContextForImage(frame)

	// load font
	fontPath := internal.GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	// lyrics file
	if _, ok := internal.InputsStore["lyrics_file"]; ok {
		sLBRS := internal.ObjCoords[internal.SelectLyricsBtn]
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(sLBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(internal.InputsStore["lyrics_file"]), 400, float64(sLBRS.OriginY)+internal.FontSize)
	}

	// font file
	if _, ok := internal.InputsStore["font_file"]; ok {
		sFFBRS := internal.ObjCoords[internal.FontFileBtn]

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(sFFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(internal.InputsStore["font_file"]), 400, float64(sFFBRS.OriginY)+internal.FontSize)
	}

	// bg file
	if _, ok := internal.InputsStore["background_file"]; ok {

		bGFBRS := internal.ObjCoords[internal.BgFileBtn]
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(bGFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(internal.InputsStore["background_file"]), 400, float64(bGFBRS.OriginY)+internal.FontSize)
	}

	// music file
	if _, ok := internal.InputsStore["music_file"]; ok {
		mFBRS := internal.ObjCoords[internal.MusicFileBtn]

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(mFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(internal.InputsStore["music_file"]), 400, float64(mFBRS.OriginY)+internal.FontSize)

	}

	// color
	if _, ok := internal.InputsStore["lyrics_color"]; ok {
		cBRS := internal.ObjCoords[internal.LyricsColorBtn]
		ggCtx.SetHexColor(internal.InputsStore["lyrics_color"])
		ggCtx.DrawRectangle(400, float64(cBRS.OriginY), 100, 40)
		ggCtx.Fill()
	}

	return ggCtx.Image()
}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range internal.ObjCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
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

		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.FontFileBtn:
		filename := pickFileUbuntu("ttf")
		if filename == "" {
			return
		}
		internal.InputsStore["font_file"] = filename
		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.BgFileBtn:
		filename := pickFileUbuntu("png")
		if filename == "" {
			return
		}
		internal.InputsStore["background_file"] = filename

		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.MusicFileBtn:
		filename := pickFileUbuntu("mp3")
		if filename == "" {
			return
		}
		internal.InputsStore["music_file"] = filename

		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.LyricsColorBtn:
		tmpColor := pickColor()
		if tmpColor == "" {
			return
		}
		internal.InputsStore["lyrics_color"] = tmpColor

		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()

	case internal.OurSite:
		internal.ExternalLaunch("https://sae.ng")

	case internal.RenderBtn:
		if len(internal.InputsStore) != 5 {
			return
		}

		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		internal.DrawRenderView(window, currentFrame)
		internal.InChannel <- "mp4"

	case internal.RenderL8fBtn:
		if len(internal.InputsStore) != 5 {
			return
		}

		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		window.SetCursorPosCallback(nil)
		internal.DrawRenderView(window, currentFrame)
		internal.InChannel <- "l8f"
	}
}

func cursorPosCB(window *glfw.Window, xpos, ypos float64) {
	if runtime.GOOS == "linux" {
		// linux fires too many events
		internal.CursorEventsCount += 1
		if internal.CursorEventsCount != 10 {
			return
		} else {
			internal.CursorEventsCount = 0
		}
	}

	wWidth, wHeight := window.GetSize()

	var widgetRS g143.RectSpecs
	var widgetCode int

	xPosInt := int(xpos)
	yPosInt := int(ypos)
	for code, RS := range internal.ObjCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {

		currentFrame := refreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()
		return
	}

	rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
		widgetRS.OriginX+widgetRS.Width,
		widgetRS.OriginY+widgetRS.Height)

	pieceOfCurrentFrame := imaging.Crop(internal.EmptyFrameNoInputs, rectA)
	invertedPiece := imaging.Invert(pieceOfCurrentFrame)

	ggCtx := gg.NewContextForImage(internal.EmptyFrameNoInputs)
	ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

	currentFrame := refreshInputsOnWindow(window, ggCtx.Image())
	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
	window.SwapBuffers()
}
