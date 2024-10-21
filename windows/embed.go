package main

import (
	_ "embed"
)

//go:embed colors.txt
var Colors2 []byte

//go:embed "ffmpeg/ffmpeg.exe"
var ffmpegBytes []byte
