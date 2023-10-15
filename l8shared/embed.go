package l8shared

import (
	_ "embed"
)

//go:embed "bmtf.txt"
var SampleLyricsFile []byte

//go:embed sae_logo.png
var SaeLogoBytes []byte

//go:embed "guitar.png"
var GuitarJPG []byte
