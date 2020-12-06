package vlcServer

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adanalife/tripbot/pkg/config"
	terrors "github.com/adanalife/tripbot/pkg/errors"
	"github.com/adanalife/tripbot/pkg/helpers"
	onscreensServer "github.com/adanalife/tripbot/pkg/onscreens-server"
	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//TODO: use more StatusExpectationFailed instead of http.StatusUnprocessableEntity
func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// healthcheck URL, for tools to verify the stream is alive
		if r.URL.Path == "/health" {
			healthCheck(w)

		} else if strings.HasPrefix(r.URL.Path, "/vlc/current") {
			// return the currently-playing file
			fmt.Fprintf(w, currentlyPlaying())

		} else if strings.HasPrefix(r.URL.Path, "/vlc/play") {
			videoFile, ok := r.URL.Query()["video"]
			if !ok || len(videoFile) > 1 {
				//TODO: eventually this could just play instead of hard-requiring a param
				http.Error(w, "417 expectation failed", http.StatusExpectationFailed)
				return
			}

			spew.Dump(videoFile)
			playVideoFile(videoFile[0])

			//TODO: better response
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/vlc/back") {
			num, ok := r.URL.Query()["n"]
			if !ok || len(num) > 1 {
				back(1)
				return
			}
			i, err := strconv.Atoi(num[0])
			if err != nil {
				terrors.Log(err, "couldn't convert input to int")
				http.Error(w, "422 unprocessable entity", http.StatusUnprocessableEntity)
				return
			}

			back(i)

			//TODO: better response
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/vlc/skip") {
			num, ok := r.URL.Query()["n"]
			if !ok || len(num) > 1 {
				skip(1)
				return
			}
			i, err := strconv.Atoi(num[0])
			if err != nil {
				terrors.Log(err, "couldn't convert input to int")
				http.Error(w, "422 unprocessable entity", http.StatusUnprocessableEntity)
				return
			}

			skip(i)

			//TODO: better response
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/vlc/random") {
			// play a random file
			err := PlayRandom()
			if err != nil {
				http.Error(w, "error playing random", http.StatusInternalServerError)
			}
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/onscreens/flag/show") {
			base64content, ok := r.URL.Query()["duration"]
			if !ok || len(base64content) > 1 {
				http.Error(w, "417 expectation failed", http.StatusExpectationFailed)
				return
			}
			durStr, err := helpers.Base64Decode(base64content[0])
			if err != nil {
				terrors.Log(err, "unable to decode string")
				http.Error(w, "422 unprocessable entity", http.StatusUnprocessableEntity)
				return
			}
			dur, err := time.ParseDuration(durStr)
			if err != nil {
				http.Error(w, "unable to parse duration", http.StatusInternalServerError)
				return
			}
			onscreensServer.ShowFlag(dur)
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/onscreens/gps/hide") {
			onscreensServer.HideGPSImage()
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/onscreens/gps/show") {
			onscreensServer.ShowGPSImage()
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/onscreens/timewarp/show") {
			onscreensServer.ShowTimewarp()
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/onscreens/leaderboard/show") {
			base64content, ok := r.URL.Query()["content"]
			if !ok || len(base64content) > 1 {
				http.Error(w, "417 expectation failed", http.StatusExpectationFailed)
				return
			}
			spew.Dump(base64content[0])
			content, err := helpers.Base64Decode(base64content[0])
			if err != nil {
				terrors.Log(err, "unable to decode string")
				http.Error(w, "422 unprocessable entity", http.StatusUnprocessableEntity)
				return
			}

			onscreensServer.Leaderboard.Show(content)
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/onscreens/middle/hide") {
			onscreensServer.MiddleText.Hide()
			fmt.Fprintf(w, "OK")

		} else if strings.HasPrefix(r.URL.Path, "/onscreens/middle/show") {
			base64content, ok := r.URL.Query()["msg"]
			if !ok || len(base64content) > 1 {
				http.Error(w, "417 expectation failed", http.StatusExpectationFailed)
				return
			}
			msg, err := helpers.Base64Decode(base64content[0])
			if err != nil {
				terrors.Log(err, "unable to decode string")
				http.Error(w, "422 unprocessable entity", http.StatusUnprocessableEntity)
				return
			}
			onscreensServer.MiddleText.Show(msg)
			fmt.Fprintf(w, "OK")

			// return a favicon if anyone asks for one
		} else if r.URL.Path == "/favicon.ico" {
			http.ServeFile(w, r, "assets/favicon.ico")

			// some other URL was used
		} else {
			http.Error(w, "404 not found", http.StatusNotFound)
			log.Println("someone tried hitting", r.URL.Path)
			return
		}

	// someone tried a PUT or a DELETE or something
	default:
		fmt.Fprintf(w, "Only GET methods are supported.\n")
	}
}

// Start starts the web server
func Start() {
	log.Println("Starting VLC web server on host", config.VlcServerHost)

	// make prometheus metrics available
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", handle)

	// ListenAndServe() wants a port in the format ":NUM"
	//TODO: error if there's no colon to split on
	port := ":" + strings.Split(config.VlcServerHost, ":")[1]
	//TODO: replace certs with autocert: https://stackoverflow.com/a/40494806
	// err := http.ListenAndServeTLS(port, "infra/tripbot.dana.lol.fullchain.pem", "infra/tripbot.dana.lol.key", nil)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		terrors.Fatal(err, "couldn't start server")
	}
}
