package main

import (
	_ "embed"
)

//go:embed "bmtf.txt"
var sampleLyricsFile []byte

//go:embed sae_logo.png
var SaeLogoBytes []byte
