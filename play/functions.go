package main

import (
	"fmt"
	"strconv"
	"strings"
)

func TimeFormatToSeconds(s string) int {
	// calculate total duration of the song
	parts := strings.Split(s, ":")
	minutesPartConverted, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	secondsPartConverted, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	totalSecondsOfSong := (60 * minutesPartConverted) + secondsPartConverted
	return totalSecondsOfSong
}

func secondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	secondsStr := fmt.Sprintf("%d", seconds)
	if seconds < 10 {
		secondsStr = "0" + secondsStr
	}
	str := fmt.Sprintf("%d:%s", minutes, secondsStr)
	return str
}
