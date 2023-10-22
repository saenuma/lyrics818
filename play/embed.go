package main

import (
	_ "embed"
)

//go:embed .asoundrc
var SoundRC []byte
