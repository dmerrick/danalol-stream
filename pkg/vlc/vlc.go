package vlc

import (
	"os"
	"path/filepath"

	theirVlc "github.com/adrg/libvlc-go"
	"github.com/davecgh/go-spew/spew"
	"github.com/dmerrick/danalol-stream/pkg/config"
	terrors "github.com/dmerrick/danalol-stream/pkg/errors"
)

var player *theirVlc.ListPlayer
var mediaList *theirVlc.MediaList

func Init() {
	var err error

	//TODO: add more flags (--no-audio?)
	if err = theirVlc.Init("--quiet"); err != nil {
		terrors.Fatal(err, "error initializing VLC")
	}

	// create a new player
	player, err = theirVlc.NewListPlayer()
	if err != nil {
		terrors.Fatal(err, "error creating VLC player")
	}
	mediaList, err = theirVlc.NewMediaList()
	if err != nil {
		terrors.Fatal(err, "error creating VLC media list")
	}

	err = player.SetMediaList(mediaList)
	if err != nil {
		terrors.Fatal(err, "error setting VLC media list")
	}

	// loop forever
	err = player.SetPlaybackMode(theirVlc.Loop)
	if err != nil {
		terrors.Fatal(err, "error setting VLC playback mode")
	}
}

func Shutdown() {
	player.Stop()
	player.Release()
	theirVlc.Release()
}

func LoadMedia() {
	var files []string

	// add all files from the VideoDir to the medialist
	err := filepath.Walk(config.VideoDir, func(path string, info os.FileInfo, err error) error {
		// skip the dir itself
		if path == config.VideoDir {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		terrors.Fatal(err, "error walking VideoDir")
	}

	for _, file := range files {
		// add the media to VLC
		err = mediaList.AddMediaFromPath(file)
		if err != nil {
			terrors.Fatal(err, "error adding files to VLC media list")
		}
	}

	spew.Dump(mediaList)
}

func Play() error {
	// start playing the media
	return player.Play()
}
