package testplay

import (
	"image"
	"time"

	g143 "github.com/bankole7782/graphics143"
)

const (
	FPS         = 24
	FontSize    = 20
	Scale       = 0.8
	MobileScale = 0.7
	fontColor   = "#444"

	SwitchViewBtn = 102
)

var (
	ObjCoords          map[int]g143.Rect
	TmpNowPlayingImg   image.Image
	currentWindowFrame image.Image
	StartTime          time.Time
	PausedSeconds      int
	CurrentPlaySeconds int

	DeviceView = "laptop"
	SongPath   string
)
