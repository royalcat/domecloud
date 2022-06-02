package entrymodel

import (
	"time"
)

type EntryInfo struct {
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"modTime"`
	MimeType MimeType

	MediaInfo *MediaInfo `json:"mediaInfo"`
}

type MediaType string

const (
	MediaTypeNone  MediaType = "none"
	MediaTypeVideo           = "video"
	MediaTypeImage           = "image"
	MediaTypeAudio           = "audio"
)

type MediaInfo struct {
	MediaType MediaType `json:"mediaType"`

	ImageInfo *ImageInfo `json:"imageInfo,omitempty"`
	VideoInfo *VideoInfo `json:"videoInfo,omitempty"`
	AudioInfo *AudioInfo `json:"audioInfo,omitempty"`
}

type ImageInfo struct {
	Resolution Resolution `json:"resolution"`
}

type VideoInfo struct {
	Duration   time.Duration `json:"duration"`
	Resolution Resolution    `json:"resolution"`
}

type AudioInfo struct {
	Duration time.Duration `json:"duration"`
}

type Resolution struct {
	Height uint64 `json:"height"`
	Width  uint64 `json:"width"`
}
