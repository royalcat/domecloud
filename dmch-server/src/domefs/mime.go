package domefs

import (
	"dmch-server/src/domefs/media"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
)

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 512

func (domefs *DomeFS) MimeType(name string) (media.MimeType, error) {
	ext := filepath.Ext(name)
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		return media.MimeType(ctype), nil
	}

	f, err := domefs.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		return "", fmt.Errorf("cant open file with error: %w", err)
	}

	var reader io.ReadSeeker = f
	var buf [sniffLen]byte
	n, _ := io.ReadFull(reader, buf[:])
	ctype = http.DetectContentType(buf[:n])
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("seeker can't seek with error: %w", err)
	}

	return media.MimeType(ctype), nil
}
