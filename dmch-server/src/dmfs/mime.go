package dmfs

import (
	"dmch-server/src/dmfs/media"
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

func (dmfs *DmFS) MimeType(name string) (media.MimeType, error) {
	stat, err := dmfs.Stat(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		return "", fmt.Errorf("cant get file stat with error: %w", err)
	}

	ctype := mime.TypeByExtension(filepath.Ext(stat.Name()))
	if ctype != "" {
		return media.MimeType(ctype), nil
	}

	f, err := dmfs.Open(name)
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
