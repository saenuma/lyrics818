package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
)

const (
	FPS             = 24
	FontSize        = 20
	OpenWDBtn       = 101
	ViewLyricsBtn   = 102
	SelectLyricsBtn = 103
	FontFileBtn     = 104
	BgFileBtn       = 105
	MusicFileBtn    = 106
	LyricsColorBtn  = 107
	RenderBtn       = 109
	OurSite         = 111
)

var (
	ObjCoords          map[int]g143.Rect
	InputsStore        map[string]string
	EmptyFrameNoInputs image.Image

	InChannel             chan bool
	InColorChannel        chan bool
	DoneWithRender        bool
	ClearAfterColorPicker bool
	PickedColor           string

	CursorEventsCount = 0
)
