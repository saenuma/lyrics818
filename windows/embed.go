package main

import (
	_ "embed"
)

//go:embed "execs/acpicker.exe"
var acPickerBytes []byte

//go:embed "execs/ffmpeg.exe"
var ffmpegBytes []byte
