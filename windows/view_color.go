package main

import (
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/internal"
)

var allColors []string

func drawPickColors(window *glfw.Window) {

	tmp := strings.ReplaceAll(string(Colors2), "\r", "")
	colors := strings.Split(tmp, "\n")
	allColors = colors

	wWidth, wHeight := window.GetSize()

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

	gutter := 10
	currentX := gutter
	currentY := gutter

	boxDimension := 50
	for i, aColor := range colors {
		ggCtx.SetHexColor(aColor)
		ggCtx.DrawRectangle(float64(currentX), float64(currentY), float64(boxDimension), float64(boxDimension))
		ggCtx.Fill()
		aColorRS := g143.RectSpecs{OriginX: currentX, OriginY: currentY, Width: boxDimension, Height: boxDimension}
		colorObjCoords[i+1] = aColorRS

		newX := currentX + boxDimension + gutter
		if newX > (wWidth - boxDimension) {
			currentY += boxDimension + gutter
			currentX = gutter
		} else {
			currentX += boxDimension + gutter
		}

	}

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}

func pickColorsMouseCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	var widgetCode int

	for code, RS := range colorObjCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	internal.InputsStore["lyrics_color"] = allColors[widgetCode-1]

	// go back
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetCursorPosCallback(internal.CursorPosCB)

	currentFrame := internal.RefreshInputsOnWindow(window, internal.EmptyFrameNoInputs)
	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, currentFrame, windowRS)
	window.SwapBuffers()
}
