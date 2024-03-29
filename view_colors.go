package main

import (
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func getColors() []string {
	colorsStr := strings.ReplaceAll(string(Colors2), "\r", "")
	colors := strings.Split(colorsStr, "\n")
	return colors
}

func drawColorDialog(window *glfw.Window) {

	colors := getColors()

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

	gutter := 5
	currentX := 20
	currentY := gutter

	boxDimension := 55
	for i, aColor := range colors {
		ggCtx.SetHexColor(aColor)
		ggCtx.DrawRoundedRectangle(float64(currentX), float64(currentY), float64(boxDimension), float64(boxDimension), 4)
		ggCtx.Fill()
		aColorRS := g143.RectSpecs{OriginX: currentX, OriginY: currentY, Width: boxDimension, Height: boxDimension}
		objCoords[i+1] = aColorRS

		newX := currentX + boxDimension + gutter
		if newX > (wWidth - boxDimension) {
			currentY += boxDimension + gutter
			currentX = 20
		} else {
			currentX += boxDimension + gutter
		}

	}

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}

func colorDialogMouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()"

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	colors := getColors()
	inputsStore["lyrics_color"] = colors[widgetCode-1]

	// load default view
	allDraws(window)
	window.SetMouseButtonCallback(mouseBtnCallback)
	displayInputs(window)
}
