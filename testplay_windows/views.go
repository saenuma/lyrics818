package main

import (
	"path/filepath"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
	"github.com/saenuma/lyrics818/internal"
	"github.com/saenuma/lyrics818/l8f"
)

func DrawNowPlayingUI(window *glfw.Window, songPath string, seconds int) {
	if DeviceView == "laptop" {
		DrawNowPlayingUILaptop(window, songPath, seconds)
	} else {
		DrawNowPlayingUIMobile(window, songPath, seconds)
	}
}

func DrawNowPlayingUILaptop(window *glfw.Window, songPath string, seconds int) {
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight)

	// draw switch device button
	theCtx.drawButtonA(SwitchViewBtn, 400, 20, "Switch View Button", "#fff", "#777")

	// Scale down the image and write frame
	currFrame, _ := l8f.ReadLaptopFrame(songPath, seconds)
	displayFrameW := int(Scale * float64((*currFrame).Bounds().Dx()))
	displayFrameH := int(Scale * float64((*currFrame).Bounds().Dy()))
	tmp := imaging.Fit(*currFrame, displayFrameW, displayFrameH, imaging.Lanczos)
	theCtx.ggCtx.DrawImage(tmp, (wWidth-displayFrameW)/2, 80)

	aStr := filepath.Base(songPath)
	aStrW, _ := theCtx.ggCtx.MeasureString(aStr)
	theCtx.ggCtx.SetHexColor("#444")

	aStrY := float64(displayFrameH) + 90 + FontSize
	theCtx.ggCtx.DrawString(aStr, (float64(wWidth)-aStrW)/2, aStrY)

	window.SetTitle(songPath + "  | L8f TestPlay")

	// write time elapsed
	elapsedTimeStr := internal.SecondsToMinutes(seconds)
	theCtx.ggCtx.DrawString(elapsedTimeStr, 50, aStrY)

	// write stop time
	totalSeconds, _ := l8f.GetVideoLength(songPath)
	stopTimeStr := internal.SecondsToMinutes(totalSeconds)
	stopTimeStrW, _ := theCtx.ggCtx.MeasureString(stopTimeStr)
	theCtx.ggCtx.DrawString(stopTimeStr, float64(wWidth)-50-stopTimeStrW, aStrY)

	// save the frame
	TmpNowPlayingImg = theCtx.ggCtx.Image()

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func DrawNowPlayingUIMobile(window *glfw.Window, songPath string, seconds int) {
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight)

	// draw switch device button
	theCtx.drawButtonA(SwitchViewBtn, 400, 20, "Switch View Button", "#fff", "#777")

	// Scale down the image and write frame
	currFrame, _ := l8f.ReadMobileFrame(songPath, seconds)
	displayFrameW := int(MobileScale * float64((*currFrame).Bounds().Dx()))
	displayFrameH := int(MobileScale * float64((*currFrame).Bounds().Dy()))
	tmp := imaging.Fit(*currFrame, displayFrameW, displayFrameH, imaging.Lanczos)
	frameX := (wWidth - displayFrameW - 300) / 2
	theCtx.ggCtx.DrawImage(tmp, frameX, 80)

	aStr := filepath.Base(songPath)
	// aStrW, _ := theCtx.ggCtx.MeasureString(aStr)
	theCtx.ggCtx.SetHexColor("#444")

	// aStrY := float64(displayFrameH) + 90 + FontSize
	textX := frameX + displayFrameW + 50
	theCtx.ggCtx.DrawString(aStr, float64(textX), 100)

	window.SetTitle(songPath + "  | L8f TestPlay")

	// write time elapsed
	elapsedTimeStr := internal.SecondsToMinutes(seconds)
	theCtx.ggCtx.DrawString(elapsedTimeStr, float64(textX), 100+40)

	// write stop time
	totalSeconds, _ := l8f.GetVideoLength(songPath)
	stopTimeStr := internal.SecondsToMinutes(totalSeconds)
	// stopTimeStrW, _ := theCtx.ggCtx.MeasureString(stopTimeStr)
	theCtx.ggCtx.DrawString(stopTimeStr, float64(textX), 100+80)

	// save the frame
	TmpNowPlayingImg = theCtx.ggCtx.Image()

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}
