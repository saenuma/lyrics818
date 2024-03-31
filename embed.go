package main

import (
	_ "embed"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed "bmtf.txt"
var SampleLyricsFile []byte
