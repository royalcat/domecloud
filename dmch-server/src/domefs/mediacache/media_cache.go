package mediacache

import (
	"context"
	"dmch-server/src/domefs/dindex"
	"dmch-server/src/domefs/media"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type MediaCache struct {
	fflock sync.Mutex

	index *dindex.DomeIndex

	cacheDir string

	log *logrus.Entry
}

func NewMediaCache(db *mongo.Database, cacheDir string) *MediaCache {
	return &MediaCache{
		index:    dindex.NewDomeIndex(db),
		cacheDir: cacheDir,
		log:      logrus.WithField("service", "mediacache"),
	}
}

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 4096

func (mw *MediaCache) GetMimeType(realpath string) (media.MimeType, error) {
	ext := filepath.Ext(realpath)
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		return media.MimeType(ctype), nil
	}

	f, err := os.Open(realpath)
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

func (mw *MediaCache) GetInfoFilePath(ctx context.Context, realpath, virtpath string) (string, error) {
	_, err := mw.genVideoInfo(ctx, realpath, virtpath)
	if err != nil {
		mw.log.Errorf("Eror generating video info: %w", err)
		return "", err
	}

	return mw.getInfoPath(virtpath), nil
}

func (mw *MediaCache) GetPreviewsDirPath(ctx context.Context, realpath, virtpath string) (string, error) {
	info, err := mw.genVideoInfo(ctx, realpath, virtpath)
	if err != nil {
		mw.log.Errorf("Eror generating video info: %w", err)
		return "", err
	}

	err = mw.genPreviews(ctx, realpath, virtpath, getTimestamps(info.VideoInfo.Duration))
	if err != nil {
		mw.log.Errorf("Eror generating video previews: %w", err)
		return "", err
	}

	return mw.getPreviewsDirPath(virtpath), nil
}
