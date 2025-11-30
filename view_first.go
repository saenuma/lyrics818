package main

import (
	"image"
	"os"
	"path/filepath"
	"runtime"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
)

func GetDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "l818_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}

func AllDraws(window *glfw.Window) {
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight)

	oWDRect := theCtx.drawButtonA(OpenWDBtn, 80, 10, "Open Working Directory", "#444", "#EAE6C7")
	vSX := nextX(oWDRect, 30)
	theCtx.drawButtonA(ViewLyricsBtn, vSX, 10, "View Sample Lyrics", "#444", "#EAE6C7")

	// Help messages
	msg1 := "1. All files must be placed in the working directory of this program."
	msg2 := "2. The background_file must be of dimensions (1366px x 768px)"

	theCtx.ggCtx.DrawString(msg1, 60, 90+FontSize)
	theCtx.ggCtx.DrawString(msg2, 60, 90+30+FontSize)

	sLBRect := theCtx.drawButtonA(SelectLyricsBtn, 40, 160, "Select Lyrics File (.txt)", "#fff", "#958E6B")
	sFBY := nextY(sLBRect, 20)
	sFBRect := theCtx.drawButtonA(FontFileBtn, 40, sFBY, "Select Font File (.ttf)", "#fff", "#958E6B")
	sBBY := nextY(sFBRect, 20)
	sBBRect := theCtx.drawButtonA(BgFileBtn, 40, sBBY, "Select Background File (.png)", "#fff", "#958E6B")
	mBBY := nextY(sBBRect, 20)
	mBBRect := theCtx.drawButtonA(MusicFileBtn, 40, mBBY, "Select Music File (.mp3)", "#fff", "#958E6B")
	lBBY := nextY(mBBRect, 20)
	lBBRect := theCtx.drawButtonA(LyricsColorBtn, 40, lBBY, "Pick Lyrics Color", "#fff", "#958E6B")

	rBX := 220
	rBY := nextY(lBBRect, 40)
	rBBrect := theCtx.drawButtonA(RenderBtn, rBX, rBY, "Make Lyrics Video (.mp4)", "#444", "#EAE6C7")
	rLBY := nextY(rBBrect, 20)
	theCtx.drawButtonA(RenderL8fBtn, rBX, rLBY, "Make Lyrics Video (.l8f)", "#444", "#EAE6C7")

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	EmptyFrameNoInputs = theCtx.ggCtx.Image()
}

func RefreshInputsOnWindow(window *glfw.Window, frame image.Image) image.Image {
	wWidth, _ := window.GetSize()

	ggCtx := gg.NewContextForImage(frame)

	// load font
	fontPath := GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	// lyrics file
	if _, ok := InputsStore["lyrics_file"]; ok {
		sLBRS := ObjCoords[SelectLyricsBtn]
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(sLBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(InputsStore["lyrics_file"]), 400, float64(sLBRS.OriginY)+FontSize)
	}

	// font file
	if _, ok := InputsStore["font_file"]; ok {
		sFFBRS := ObjCoords[FontFileBtn]

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(sFFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(InputsStore["font_file"]), 400, float64(sFFBRS.OriginY)+FontSize)
	}

	// bg file
	if _, ok := InputsStore["background_file"]; ok {

		bGFBRS := ObjCoords[BgFileBtn]
		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(bGFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(InputsStore["background_file"]), 400, float64(bGFBRS.OriginY)+FontSize)
	}

	// music file
	if _, ok := InputsStore["music_file"]; ok {
		mFBRS := ObjCoords[MusicFileBtn]

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawRectangle(400, float64(mFBRS.OriginY), float64(wWidth)-400, 40)
		ggCtx.Fill()

		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(filepath.Base(InputsStore["music_file"]), 400, float64(mFBRS.OriginY)+FontSize)

	}

	// color
	if _, ok := InputsStore["lyrics_color"]; ok {
		cBRS := ObjCoords[LyricsColorBtn]
		ggCtx.SetHexColor(InputsStore["lyrics_color"])
		ggCtx.DrawRectangle(400, float64(cBRS.OriginY), 100, 40)
		ggCtx.Fill()
	}

	return ggCtx.Image()
}

func CursorPosCB(window *glfw.Window, xpos, ypos float64) {
	if runtime.GOOS == "linux" {
		// linux fires too many events
		CursorEventsCount += 1
		if CursorEventsCount != 10 {
			return
		} else {
			CursorEventsCount = 0
		}
	}

	wWidth, wHeight := window.GetSize()

	var widgetRS g143.Rect
	var widgetCode int

	xPosInt := int(xpos)
	yPosInt := int(ypos)
	for code, RS := range ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {

		currentFrame := RefreshInputsOnWindow(window, EmptyFrameNoInputs)
		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
		window.SwapBuffers()
		return
	}

	rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
		widgetRS.OriginX+widgetRS.Width,
		widgetRS.OriginY+widgetRS.Height)

	pieceOfCurrentFrame := imaging.Crop(EmptyFrameNoInputs, rectA)
	invertedPiece := imaging.AdjustBrightness(pieceOfCurrentFrame, -20)

	ggCtx := gg.NewContextForImage(EmptyFrameNoInputs)
	ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

	currentFrame := RefreshInputsOnWindow(window, ggCtx.Image())
	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
	window.SwapBuffers()
}
