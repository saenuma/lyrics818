package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"text/template"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/saenuma/lyrics818/l8f"
	"github.com/saenuma/lyrics818/l8shared"
	sDialog "github.com/sqweek/dialog"
)

func startBackend() {
	rootPath, err := l8shared.GetRootPath()
	if err != nil {
		panic(err)
	}
	playerPath := filepath.Join(rootPath, ".player")
	os.MkdirAll(playerPath, 0777)

	r := mux.NewRouter()

	tmpAudioPath := filepath.Join(playerPath, "tmp_audio.mp3")

	currentVideoPath := ""
	currentDevice := ""

	r.HandleFunc("/gs/{obj}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rawObj, err := contentStatics.ReadFile("statics/" + vars["obj"])
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+vars["obj"])
		contentType := http.DetectContentType(rawObj)
		w.Header().Set("Content-Type", contentType)
		w.Write(rawObj)
	})

	r.HandleFunc("/xdg/", func(w http.ResponseWriter, r *http.Request) {
		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/C", "start", r.FormValue("p")).Run()
		} else if runtime.GOOS == "linux" {
			exec.Command("xdg-open", r.FormValue("p")).Run()
		}
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		type Context struct {
			OutPath string
		}
		tmpl := template.Must(template.ParseFS(content, "templates/app.html"))
		tmpl.Execute(w, Context{rootPath})
	})

	r.HandleFunc("/pick_l818_file", func(w http.ResponseWriter, r *http.Request) {
		filename, err := sDialog.File().Filter("lyrics818 video", "l8f").Load()
		if err == nil {
			fmt.Fprint(w, filename)
		}
	})

	r.HandleFunc("/begin_player", func(w http.ResponseWriter, r *http.Request) {
		currentVideoPath = r.FormValue("vid_file")

		audioBytes, err := l8f.ReadAudio(r.FormValue("vid_file"))
		if err != nil {
			panic(err)
		}

		currentDevice = r.FormValue("device")
		os.WriteFile(tmpAudioPath, audioBytes, 0777)
		fmt.Println("ok")
	})

	r.HandleFunc("/get_audio", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, tmpAudioPath)
	})

	r.HandleFunc("/get_frame/{number}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		frameNumStr := vars["number"]

		frameNum, _ := strconv.Atoi(frameNumStr)
		if currentDevice == "laptop" {
			currFrame, _ := l8f.ReadLaptopFrame(currentVideoPath, frameNum)
			imaging.Save(*currFrame, filepath.Join(playerPath, "frame.png"))
		} else {
			currFrame, _ := l8f.ReadMobileFrame(currentVideoPath, frameNum)
			imaging.Save(*currFrame, filepath.Join(playerPath, "frame.png"))

		}

		http.ServeFile(w, r, filepath.Join(playerPath, "frame.png"))
	})

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
