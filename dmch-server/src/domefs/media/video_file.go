package media

import (
	"time"
)

type Resolution struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type VideoInfo struct {
	Path       string        `json:"path"`
	Size       int64         `json:"size"`
	ModTime    time.Time     `json:"modTime"`
	Duration   time.Duration `json:"duration"`
	Resolution Resolution    `json:"resolution"`
}
