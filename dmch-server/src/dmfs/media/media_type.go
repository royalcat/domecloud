package media

import (
	"strings"
)

type MimeType string

func (mimetype MimeType) MediaType() MediaType {
	mtlower := strings.ToLower(string(mimetype))
	if strings.HasPrefix(mtlower, MediaTypeVideo) {
		return MediaTypeVideo
	}
	if strings.HasPrefix(mtlower, MediaTypeImage) {
		return MediaTypeVideo
	}

	return MediaTypeNone
}

type MediaType string

const (
	MediaTypeNone  MediaType = "none"
	MediaTypeVideo           = "video"
	MediaTypeImage           = "image"
)

func endsWithOneOf(s string, ends []string) bool {
	for _, end := range ends {
		if strings.HasSuffix(s, end) {
			return true
		}
	}
	return false
}
