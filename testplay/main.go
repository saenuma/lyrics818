package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jchv/go-webview2"
)

const port = "31992"

func main() {
	debug := false
	if os.Getenv("SAENUMA_DEVELOPER") == "true" {
		debug = true
	}

	go startBackend()

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     debug,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title: "testplay - Test Videos made with Lyrics818",
		},
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()
	w.SetSize(1200, 600, webview2.HintNone)
	w.Navigate(fmt.Sprintf("http://127.0.0.1:%s", port))
	w.Run()
}
