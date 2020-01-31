package video

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/dmerrick/danalol-stream/pkg/background"
	terrors "github.com/dmerrick/danalol-stream/pkg/errors"
	"github.com/dmerrick/danalol-stream/pkg/helpers"
	"github.com/logrusorgru/aurora"
)

// CurrentlyPlaying is the video that is currently playing
var CurrentlyPlaying Video

// these are used to keep track of the current video
var curVid, preVid string
var timeStarted time.Time

// GetCurrentlyPlaying will use lsof to figure out
// which dashcam video is currently playing (seriously)
//TODO: consider making this return a video struct
func GetCurrentlyPlaying() {
	var err error
	// save the video we used last time
	preVid = curVid
	// figure out whats currently playing
	curVid = figureOutCurrentVideo()

	// if the currently-playing video has changed
	if curVid != preVid {
		log.Printf("now playing %s", aurora.Yellow(curVid))

		// reset the stopwatch
		timeStarted = time.Now()

		// share the Video with the system
		CurrentlyPlaying, err = LoadOrCreate(curVid)
		if err != nil {
			terrors.Log(err, fmt.Sprintf("unable to create Video from %s", curVid))
		}

		// show the no-GPS image
		if CurrentlyPlaying.Flagged {
			log.Println("SHOWING NO-GPS IMAGE")
			background.GPSImage.Show("", 60*time.Second)
		} else {
			log.Println("HIDING NO-GPS IMAGE")
			background.GPSImage.Hide()
		}
	}
}

// CurrentProgress represents how long the video has been playing
// it will be useful eventually for choosing the exact right screenshot
func CurrentProgress() time.Duration {
	return time.Since(timeStarted)
}

func figureOutCurrentVideo() string {
	// run the shell script to get currently-playing video
	scriptPath := path.Join(helpers.ProjectRoot(), "bin/current-file.sh")
	out, err := exec.Command(scriptPath).Output()
	outString := strings.TrimSpace(string(out))
	if err != nil {
		terrors.Log(err, "failed to get currently-playing video")
		log.Println(outString)
		return ""
	}
	return outString
}
