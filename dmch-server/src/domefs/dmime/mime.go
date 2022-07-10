package dmime

import "strings"

type MimeType string

const (
	MimeTypeJson      MimeType = "application/json"
	MimeTypeDirectory MimeType = "inode/directory"
)

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

func endsWithOneOf(s string, ends []string) bool {
	for _, end := range ends {
		if strings.HasSuffix(s, end) {
			return true
		}
	}
	return false
}
