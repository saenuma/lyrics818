package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/saenuma/lyrics818/internal"
)

type Ctx struct {
	WindowWidth  int
	WindowHeight int
	ggCtx        *gg.Context
	ObjCoords    *map[int]g143.Rect
}

func New2dCtx(wWidth, wHeight int) Ctx {
	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := internal.GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: wWidth, WindowHeight: wHeight, ggCtx: ggCtx}
	return ctx
}

func Continue2dCtx(img image.Image, objCoords *map[int]g143.Rect) Ctx {
	ggCtx := gg.NewContextForImage(img)

	// load font
	fontPath := internal.GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: img.Bounds().Dx(), WindowHeight: img.Bounds().Dy(), ggCtx: ggCtx,
		ObjCoords: objCoords}
	return ctx
}

func (ctx *Ctx) drawButtonA(btnId, originX, originY int, text, textColor, bgColor string) g143.Rect {
	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+20, textH+15
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+10, float64(originY)+FontSize)

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	ObjCoords[btnId] = btnARect
	return btnARect
}

func (ctx *Ctx) windowRect() g143.Rect {
	return g143.NewRect(0, 0, ctx.WindowWidth, ctx.WindowHeight)
}

// func nextHorizontalCoords(aRect g143.Rect, margin int) (int, int) {
// 	nextOriginX := aRect.OriginX + aRect.Width + margin
// 	nextOriginY := aRect.OriginY
// 	return nextOriginX, nextOriginY
// }

// func nextVerticalCoords(aRect g143.Rect, margin int) (int, int) {
// 	nextOriginX := margin
// 	nextOriginY := aRect.OriginY + aRect.Height + margin
// 	return nextOriginX, nextOriginY
// }
