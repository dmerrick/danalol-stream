package video

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"time"

	c "github.com/adanalife/tripbot/pkg/config/tripbot"
)

// Videos represent a video file containing dashcam footage
type Video struct {
	Id          int           `db:"id"`
	Slug        string        `db:"slug"`
	Lat         float64       `db:"lat"`
	Lng         float64       `db:"lng"`
	NextVid     sql.NullInt64 `db:"next_vid"`
	PrevVid     sql.NullInt64 `db:"prev_vid"`
	Flagged     bool          `db:"flagged"`
	State       string        `db:"state"`
	DateFilmed  time.Time     `db:"date_filmed"`
	DateCreated time.Time     `db:"date_created"`
}

// Location returns a lat/lng pair
//TODO: refactor out the error return value
func (v Video) Location() (float64, float64, error) {
	var err error
	if v.Flagged {
		err = errors.New("video is flagged")
	}
	return v.Lat, v.Lng, err
}

//TODO: add color, include location/state/lat/lng?
//TODO: where else does this get used tho?
// ex: 2018_0514_224801_013_a_opt
func (v Video) String() string {
	return v.Slug
}

// a DashStr is the string we get from the dashcam
// an example file: 2018_0514_224801_013.MP4
// an example dashstr: 2018_0514_224801_013
// ex: 2018_0514_224801_013
func (v Video) DashStr() string {
	//TODO: this never should have happened, but it did and it crashed the bot
	if len(v.Slug) < 20 {
		return ""
	}
	return v.Slug[:20]
}

// ex: 2018_0514_224801_013.MP4
func (v Video) File() string {
	return fmt.Sprintf("%s.MP4", v.Slug)
}

// ex: /Volumes/.../2018_0514_224801_013.MP4
func (v Video) Path() string {
	return filepath.Join(c.Conf.VideoDir, v.File())
}

// toDate parses the vidStr and returns a time.Time object for the video
func (v Video) toDate() time.Time {
	vidStr := v.String()
	year, _ := strconv.Atoi(vidStr[:4])
	month, _ := strconv.Atoi(vidStr[5:7])
	day, _ := strconv.Atoi(vidStr[7:9])
	hour, _ := strconv.Atoi(vidStr[10:12])
	minute, _ := strconv.Atoi(vidStr[12:14])
	second, _ := strconv.Atoi(vidStr[14:16])

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	return t
}

// slug strips the path and extension off the file
func slug(file string) string {
	fileName := path.Base(file)
	return removeFileExtension(fileName)
}

func removeFileExtension(filename string) string {
	ext := path.Ext(filename)
	return filename[0 : len(filename)-len(ext)]
}
