package main

import (
	_ "embed"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed "bmtf.txt"
var SampleLyricsFile []byte

//go:embed colors.txt
var Colors2 []byte

// // go:embed "ffmpeg/ffmpeg.exe"
var ffmpegBytes []byte
