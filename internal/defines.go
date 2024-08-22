package internal

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
	RenderL8fBtn    = 110
	OurSite         = 111
)

var EmptyFrameNoInputs image.Image

var InputsStore map[string]string

var InChannel chan string
var ClearAfterRender bool

var CursorEventsCount = 0

var ObjCoords map[int]g143.RectSpecs
