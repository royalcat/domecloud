package media

import (
	"time"
)

type Resolution struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type VideoInfo struct {
	Duration time.Duration `json:"duration"`
}

type VisualMediaInfo struct {
	Path       string     `json:"path"`
	Size       int64      `json:"size"`
	ModTime    time.Time  `json:"modTime"`
	Resolution Resolution `json:"resolution"`
	MediaType  MimeType

	VideoInfo VideoInfo `json:"videoInfo"`
}
