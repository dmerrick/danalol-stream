package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/dmerrick/danalol-stream/pkg/config"
	"github.com/dmerrick/danalol-stream/pkg/database"
	terrors "github.com/dmerrick/danalol-stream/pkg/errors"
	"github.com/dmerrick/danalol-stream/pkg/helpers"
	"github.com/dmerrick/danalol-stream/pkg/video"
	"github.com/joho/godotenv"
	"github.com/kelvins/geocoder"
)

// this will hold the filename passed in via the CLI
var videoFile string
var current bool

//TODO: this isn't a videoFile anymore
func process(videoFile string) {
	log.Println("working on:", videoFile)

	vid, err := video.LoadOrCreate(videoFile)
	if err != nil {
		terrors.Log(err, fmt.Sprintf("failed to create video: %v", videoFile))
		return
	}

	lat, lon, err := vid.Location()
	if err != nil {
		terrors.Log(err, fmt.Sprintf("failed to process image: %v", videoFile))
		return
	}
	url := helpers.GoogleMapsURL(lat, lon)
	fmt.Println(url)
}

func screencapDir() string {
	// index 11 corresponds to 245 (aka 2m45s)
	// which should have the least screencaps
	return path.Join(config.ScreencapDir(), config.TimestampsToTry[11])
}

func main() {

	log.Println("going to loop over:", screencapDir())

	// make sure we were given enough args
	checkArgs()

	// a file was passed in via the CLI
	if videoFile != "" {
		process(videoFile)

	} else {
		// loop over every file in the screencapDir
		pathToWalk := screencapDir()
		err := filepath.Walk(pathToWalk,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// skip the directory name itself
				if path == pathToWalk {
					return nil
				}

				// actually process the image
				process(path)

				return err
			})
		if err != nil {
			terrors.Log(err, "something went wrong walking the directory")
		}
	}

}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("DASHCAM_DIR") == "" {
		panic("You must set DASHCAM_DIR")
	}
	googleMapsAPIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if googleMapsAPIKey == "" {
		panic("You must set GOOGLE_MAPS_API_KEY")
	}
	geocoder.ApiKey = googleMapsAPIKey

	flag.StringVar(&videoFile, "file", "", "File to load")
	flag.BoolVar(&current, "current", false, "Use currently-playing video")
	flag.Parse()

	// initialize the SQL database
	database.DBCon, err = database.Initialize()
	if err != nil {
		log.Fatal("error initializing the DB", err)
	}
}

func checkArgs() {
	// set videoFile if -current was passed in
	if current {
		// first we check if too many flags were used
		if videoFile != "" {
			log.Fatal("you cannot use -current and -file at the same time")
		}
		// preload the currently-playing vid
		video.GetCurrentlyPlaying()
		videoFile = video.CurrentlyPlaying.String()
	}
}