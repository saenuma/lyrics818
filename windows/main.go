package main

import (
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sqweek/dialog"
)

const (
	fps             = 10
	fontSize        = 20
	OpenWDBtn       = 101
	ViewLyricsBtn   = 102
	SelectLyricsBtn = 103
	FontFileBtn     = 104
	BgFileBtn       = 105
	MusicFileBtn    = 106
	LyricsColorBtn  = 107
	RenderBtn       = 109
	RenderL8fBtn    = 110
	OurSite         = 111
)

var objCoords map[int]g143.RectSpecs
var currentWindowFrame image.Image
var inputsStore map[string]string

var inChannel chan string
var clearAfterRender bool

func main() {
	// _, err := v3shared.GetRootPath()
	// if err != nil {
	// 	panic(err)
	// }

	runtime.LockOSThread()

	objCoords = make(map[int]g143.RectSpecs)
	inputsStore = make(map[string]string)
	inChannel = make(chan string)

	window := g143.NewWindow(1000, 800, "lyrics818: a more comfortable lyrics video generator", false)
	drawDefaultUI(window)

	go func() {
		for {
			method := <-inChannel
			if method == "mp4" {
				ffPath := GetFFMPEGCommand()
				_, err := MakeVideo(inputsStore, ffPath)
				if err != nil {
					log.Println(err)
					return
				}

			} else if method == "l8f" {
				_, err := MakeVideoL8F(inputsStore)
				if err != nil {
					log.Println(err)
					return
				}
			}
			clearAfterRender = true
		}
	}()

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// respond to the keyboard
	// window.SetKeyCallback(keyCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		if clearAfterRender {
			// clear the UI and redraw
			inputsStore = make(map[string]string)
			drawDefaultUI(window)
			drawEndRenderView(window, currentWindowFrame)
			time.Sleep(5 * time.Second)
			drawDefaultUI(window)
			// register the ViewMain mouse callback
			window.SetMouseButtonCallback(mouseBtnCallback)
			clearAfterRender = false
		}

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}

}

func getDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "l818_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}

func drawDefaultUI(window *glfw.Window) {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	// open working directory button
	beginXOffset := 200
	ggCtx.SetHexColor("#D09090")
	owdStr := "Open Working Directory"
	owdStrW, owdStrH := ggCtx.MeasureString(owdStr)
	ggCtx.DrawRoundedRectangle(float64(beginXOffset), 10, owdStrW+50, owdStrH+25, (owdStrH+25)/2)
	ggCtx.Fill()

	owdBtnRS := g143.RectSpecs{Width: int(owdStrW) + 50, Height: int(owdStrH) + 25, OriginX: beginXOffset, OriginY: 10}
	objCoords[OpenWDBtn] = owdBtnRS

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(owdStr, float64(beginXOffset)+25, 35)

	// view sample lyrics button
	ggCtx.SetHexColor("#90D092")
	vslStr := "View Sample Lyrics"
	vslStrWidth, vslStrHeight := ggCtx.MeasureString(vslStr)
	nexBtnOriginX := owdBtnRS.OriginX + owdBtnRS.Width + 30
	ggCtx.DrawRoundedRectangle(float64(nexBtnOriginX), 10, vslStrWidth+50, vslStrHeight+25, (vslStrHeight+25)/2)
	ggCtx.Fill()

	vslBtnRS := g143.RectSpecs{Width: int(vslStrWidth) + 50, Height: int(vslStrHeight) + 25, OriginX: nexBtnOriginX,
		OriginY: 10}
	objCoords[ViewLyricsBtn] = vslBtnRS

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(vslStr, float64(vslBtnRS.OriginX)+25, 35)

	// Help messages
	ggCtx.LoadFontFace(fontPath, 30)
	ggCtx.DrawString("Help", 40, 50+30)
	ggCtx.LoadFontFace(fontPath, 20)

	msg1 := "1. All files must be placed in the working directory of this program."
	msg2 := "2. The background_file must be of dimensions (1366px x 768px)"

	ggCtx.DrawString(msg1, 60, 90+fontSize)
	ggCtx.DrawString(msg2, 60, 90+30+fontSize)

	// lyrics file button
	lfStr := "Select Lyrics File (.txt)"
	lfStrW, _ := ggCtx.MeasureString(lfStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRoundedRectangle(40, 160, lfStrW+40, 40, 10)
	ggCtx.Fill()

	lfrs := g143.NRectSpecs(40, 160, int(lfStrW+40), 40)
	objCoords[SelectLyricsBtn] = lfrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(lfStr, 60, 165+fontSize)

	// font file button
	ffStr := "Select Font File (.ttf)"
	ffStrW, _ := ggCtx.MeasureString(ffStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRoundedRectangle(40, 220, ffStrW+40, 40, 10)
	ggCtx.Fill()

	ffrs := g143.NRectSpecs(40, 220, int(ffStrW+40), 40)
	objCoords[FontFileBtn] = ffrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(ffStr, 60, 225+fontSize)

	// background file button
	bfStr := "Select Background File (.png)"
	bfStrW, _ := ggCtx.MeasureString(bfStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRoundedRectangle(40, 280, bfStrW+40, 40, 10)
	ggCtx.Fill()

	bfrs := g143.NRectSpecs(40, 280, int(bfStrW+40), 40)
	objCoords[BgFileBtn] = bfrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(bfStr, 60, 285+fontSize)

	// music file button
	mfStr := "Select Music File (.mp3)"
	mfStrW, _ := ggCtx.MeasureString(mfStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRoundedRectangle(40, 340, mfStrW+40, 40, 10)
	ggCtx.Fill()

	mfrs := g143.NRectSpecs(40, 340, int(mfStrW+40), 40)
	objCoords[MusicFileBtn] = mfrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(mfStr, 60, 345+fontSize)

	// lyrics color button
	lcStr := "Pick Lyrics Color"
	lcStrW, _ := ggCtx.MeasureString(lcStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRoundedRectangle(40, 400, lcStrW+40, 40, 10)
	ggCtx.Fill()

	lcrs := g143.NRectSpecs(40, 400, int(lcStrW+40), 40)
	objCoords[LyricsColorBtn] = lcrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(lcStr, 60, 405+fontSize)

	// render button
	beginXOffset2 := 250
	ggCtx.SetHexColor("#A965B5")
	rStr := "Make Lyrics Video (.mp4)"
	rStrW, rStrH := ggCtx.MeasureString(rStr)
	ggCtx.DrawRoundedRectangle(float64(beginXOffset2), 480, rStrW+50, rStrH+25, (rStrH+25)/2)
	ggCtx.Fill()

	rBtnRS := g143.RectSpecs{Width: int(rStrW) + 50, Height: int(rStrH) + 25, OriginX: beginXOffset2, OriginY: 480}
	objCoords[RenderBtn] = rBtnRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(rStr, float64(beginXOffset2)+25, 485+fontSize)

	// render l8f button
	rl8X := beginXOffset2 + rBtnRS.Width + 50
	ggCtx.SetHexColor("#674C6A")
	rl8L := "Make Lyrics Video (.l8f)"
	rl8LW, rl8LH := ggCtx.MeasureString(rl8L)
	ggCtx.DrawRoundedRectangle(float64(rl8X), 480, rl8LW+50, rl8LH+25, (rl8LH+25)/2)
	ggCtx.Fill()

	rl8BtnRS := g143.NRectSpecs(rl8X, 480, int(rl8LW)+50, int(rl8LH)+25)
	objCoords[RenderL8fBtn] = rl8BtnRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(rl8L, float64(rl8X)+25, 485+fontSize)

	// draw our site below
	ggCtx.SetHexColor("#9C5858")
	fromAddr := "sae.ng"
	fromAddrWidth, fromAddrHeight := ggCtx.MeasureString(fromAddr)
	fromAddrOriginX := (wWidth - int(fromAddrWidth)) / 2
	ggCtx.DrawString(fromAddr, float64(fromAddrOriginX), float64(wHeight-int(fromAddrHeight)))
	fars := g143.RectSpecs{OriginX: fromAddrOriginX, OriginY: wHeight - 40,
		Width: int(fromAddrWidth), Height: 40}
	objCoords[OurSite] = fars

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
}

func refreshInputsOnWindow(window *glfw.Window) {
	wWidth, wHeight := window.GetSize()

	ggCtx := gg.NewContextForImage(currentWindowFrame)

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	// lyrics file
	if _, ok := inputsStore["lyrics_file"]; ok {
		sLBRS := objCoords[SelectLyricsBtn]
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(sLBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(inputsStore["lyrics_file"]), 400, float64(sLBRS.OriginY)+fontSize)
	}

	// font file
	if _, ok := inputsStore["font_file"]; ok {
		sFFBRS := objCoords[FontFileBtn]

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(sFFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(inputsStore["font_file"]), 400, float64(sFFBRS.OriginY)+fontSize)
	}

	// bg file
	if _, ok := inputsStore["background_file"]; ok {

		bGFBRS := objCoords[BgFileBtn]
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(bGFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(inputsStore["background_file"]), 400, float64(bGFBRS.OriginY)+fontSize)
	}

	// music file
	if _, ok := inputsStore["music_file"]; ok {
		mFBRS := objCoords[MusicFileBtn]

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(mFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(inputsStore["music_file"]), 400, float64(mFBRS.OriginY)+fontSize)

	}

	// color
	if _, ok := inputsStore["lyrics_color"]; ok {
		cBRS := objCoords[LyricsColorBtn]
		ggCtx.SetHexColor(inputsStore["lyrics_color"])
		ggCtx.DrawRectangle(400, float64(cBRS.OriginY), 100, 40)
		ggCtx.Fill()
	}

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	switch widgetCode {
	case OpenWDBtn:
		rootPath, _ := GetRootPath()
		externalLaunch(rootPath)

	case ViewLyricsBtn:
		drawSampleLyricsDialog(window, currentWindowFrame)

	case DialogCloseButton:
		drawDefaultUI(window)

	case SelectLyricsBtn:
		filename, err := dialog.File().Filter("Lyrics File", "txt").Load()
		if filename == "" || err != nil {
			return
		}

		inputsStore["lyrics_file"] = filename

		// write lyrics file
		drawDefaultUI(window)
		refreshInputsOnWindow(window)

	case FontFileBtn:
		filename, err := dialog.File().Filter("Font file", "ttf").Load()
		if filename == "" || err != nil {
			return
		}

		inputsStore["font_file"] = filename

		drawDefaultUI(window)
		refreshInputsOnWindow(window)

	case BgFileBtn:
		filename, err := dialog.File().Filter("PNG Image", "png").Load()
		if filename == "" || err != nil {
			return
		}

		inputsStore["background_file"] = filename

		drawDefaultUI(window)
		refreshInputsOnWindow(window)

	case MusicFileBtn:
		filename, err := dialog.File().Filter("MP3 Audio", "mp3").Load()
		if filename == "" || err != nil {
			return
		}

		inputsStore["music_file"] = filename
		drawDefaultUI(window)
		refreshInputsOnWindow(window)

	case LyricsColorBtn:

		drawPickColors(window)
		window.SetMouseButtonCallback(pickColorsMouseCallback)

		// tmpColor := pickColor()
		// if tmpColor == "" {
		// 	return
		// }
		// inputsStore["lyrics_color"] = tmpColor

		// // show sample color
		// ggCtx := gg.NewContextForImage(currentWindowFrame)

		// ggCtx.SetHexColor(tmpColor)
		// ggCtx.DrawRectangle(400, float64(widgetRS.OriginY), 100, 40)
		// ggCtx.Fill()

		// // send the frame to glfw window
		// windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		// g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		// window.SwapBuffers()

		// // save the frame
		// currentWindowFrame = ggCtx.Image()

	case OurSite:

		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "start", "https://sae.ng").Run()
		} else if runtime.GOOS == "linux" {
			exec.Command("xdg-open", "https://sae.ng").Run()
		}

	case RenderBtn:
		if len(inputsStore) != 5 {
			return
		}

		drawRenderView(window, currentWindowFrame)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		inChannel <- "mp4"

	case RenderL8fBtn:
		if len(inputsStore) != 5 {
			return
		}

		drawRenderView(window, currentWindowFrame)
		window.SetMouseButtonCallback(nil)
		window.SetKeyCallback(nil)
		inChannel <- "l8f"
	}

}
