package lyrics818

import (
	"image"
	"os"
	"path/filepath"
	"runtime"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/internal"
)

func GetDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "l818_font.ttf")
	os.WriteFile(fontPath, internal.DefaultFont, 0777)
	return fontPath
}

func AllDraws(window *glfw.Window) {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	// open working directory button
	beginXOffset := 200
	ggCtx.SetHexColor("#D09090")
	owdStr := "Open Working Directory"
	owdStrW, owdStrH := ggCtx.MeasureString(owdStr)
	ggCtx.DrawRectangle(float64(beginXOffset), 10, owdStrW+50, owdStrH+25)
	ggCtx.Fill()

	owdBtnRS := g143.Rect{Width: int(owdStrW) + 50, Height: int(owdStrH) + 25, OriginX: beginXOffset, OriginY: 10}
	ObjCoords[OpenWDBtn] = owdBtnRS

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(owdStr, float64(beginXOffset)+25, 35)

	// view sample lyrics button
	ggCtx.SetHexColor("#90D092")
	vslStr := "View Sample Lyrics"
	vslStrWidth, vslStrHeight := ggCtx.MeasureString(vslStr)
	nexBtnOriginX := owdBtnRS.OriginX + owdBtnRS.Width + 30
	ggCtx.DrawRectangle(float64(nexBtnOriginX), 10, vslStrWidth+50, vslStrHeight+25)
	ggCtx.Fill()

	vslBtnRS := g143.Rect{Width: int(vslStrWidth) + 50, Height: int(vslStrHeight) + 25, OriginX: nexBtnOriginX,
		OriginY: 10}
	ObjCoords[ViewLyricsBtn] = vslBtnRS

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(vslStr, float64(vslBtnRS.OriginX)+25, 35)

	// Help messages
	ggCtx.LoadFontFace(fontPath, 30)
	ggCtx.DrawString("Help", 40, 50+30)
	ggCtx.LoadFontFace(fontPath, 20)

	msg1 := "1. All files must be placed in the working directory of this program."
	msg2 := "2. The background_file must be of dimensions (1366px x 768px)"

	ggCtx.DrawString(msg1, 60, 90+FontSize)
	ggCtx.DrawString(msg2, 60, 90+30+FontSize)

	// lyrics file button
	lfStr := "Select Lyrics File (.txt)"
	lfStrW, _ := ggCtx.MeasureString(lfStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRectangle(40, 160, lfStrW+40, 40)
	ggCtx.Fill()

	lfrs := g143.NewRect(40, 160, int(lfStrW+40), 40)
	ObjCoords[SelectLyricsBtn] = lfrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(lfStr, 60, 165+FontSize)

	// font file button
	ffStr := "Select Font File (.ttf)"
	ffStrW, _ := ggCtx.MeasureString(ffStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRectangle(40, 220, ffStrW+40, 40)
	ggCtx.Fill()

	ffrs := g143.NewRect(40, 220, int(ffStrW+40), 40)
	ObjCoords[FontFileBtn] = ffrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(ffStr, 60, 225+FontSize)

	// background file button
	bfStr := "Select Background File (.png)"
	bfStrW, _ := ggCtx.MeasureString(bfStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRectangle(40, 280, bfStrW+40, 40)
	ggCtx.Fill()

	bfrs := g143.NewRect(40, 280, int(bfStrW+40), 40)
	ObjCoords[BgFileBtn] = bfrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(bfStr, 60, 285+FontSize)

	// music file button
	mfStr := "Select Music File (.mp3)"
	mfStrW, _ := ggCtx.MeasureString(mfStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRectangle(40, 340, mfStrW+40, 40)
	ggCtx.Fill()

	mfrs := g143.NewRect(40, 340, int(mfStrW+40), 40)
	ObjCoords[MusicFileBtn] = mfrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(mfStr, 60, 345+FontSize)

	// lyrics color button
	lcStr := "Pick Lyrics Color"
	lcStrW, _ := ggCtx.MeasureString(lcStr)
	ggCtx.SetHexColor("#5F699F")
	ggCtx.DrawRectangle(40, 400, lcStrW+40, 40)
	ggCtx.Fill()

	lcrs := g143.NewRect(40, 400, int(lcStrW+40), 40)
	ObjCoords[LyricsColorBtn] = lcrs

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(lcStr, 60, 405+FontSize)

	// render button
	beginXOffset2 := 200
	ggCtx.SetHexColor("#A965B5")
	rStr := "Make Lyrics Video (.mp4)"
	rStrW, rStrH := ggCtx.MeasureString(rStr)
	ggCtx.DrawRectangle(float64(beginXOffset2), 500, rStrW+70, rStrH+25)
	ggCtx.Fill()
	ggCtx.SetHexColor("#5D435E")
	ggCtx.DrawRoundedRectangle(float64(beginXOffset2)+rStrW+40, 500+10, 20, 20, 10)
	ggCtx.Fill()

	rBtnRS := g143.Rect{Width: int(rStrW) + 70, Height: int(rStrH) + 25, OriginX: beginXOffset2, OriginY: 500}
	ObjCoords[RenderBtn] = rBtnRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(rStr, float64(beginXOffset2)+25, 505+FontSize)

	// render l8f button
	rl8X := beginXOffset2 + rBtnRS.Width + 50
	ggCtx.SetHexColor("#674C6A")
	rl8L := "Make Lyrics Video (.l8f)"
	rl8LW, rl8LH := ggCtx.MeasureString(rl8L)
	ggCtx.DrawRectangle(float64(rl8X), 500, rl8LW+70, rl8LH+25)
	ggCtx.Fill()

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawRoundedRectangle(float64(rl8X)+rl8LW+40, 500+10, 20, 20, 10)
	ggCtx.Fill()

	rl8BtnRS := g143.NewRect(rl8X, 500, int(rl8LW)+70, int(rl8LH)+25)
	ObjCoords[RenderL8fBtn] = rl8BtnRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(rl8L, float64(rl8X)+25, 505+FontSize)

	// draw our site below
	ggCtx.SetHexColor("#9C5858")
	fromAddr := "sae.ng"
	fromAddrWidth, fromAddrHeight := ggCtx.MeasureString(fromAddr)
	fromAddrOriginX := (wWidth - int(fromAddrWidth)) / 2
	ggCtx.DrawString(fromAddr, float64(fromAddrOriginX), float64(wHeight-int(fromAddrHeight)))
	fars := g143.Rect{OriginX: fromAddrOriginX, OriginY: wHeight - 40,
		Width: int(fromAddrWidth), Height: 40}
	ObjCoords[OurSite] = fars

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	EmptyFrameNoInputs = ggCtx.Image()
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
	invertedPiece := imaging.Invert(pieceOfCurrentFrame)

	ggCtx := gg.NewContextForImage(EmptyFrameNoInputs)
	ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

	currentFrame := RefreshInputsOnWindow(window, ggCtx.Image())
	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
	window.SwapBuffers()
}
