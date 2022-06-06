package entrycache

import (
	"context"
	"dmch-server/src/domefs/dindex"
	"dmch-server/src/domefs/domefile"
	"dmch-server/src/domefs/entrymodel"
	"dmch-server/src/domefs/pathctx"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

type EntryIndex struct {
	fflock sync.Mutex

	index *dindex.DomeIndex

	pathCtx *pathctx.PathCtx

	log *logrus.Entry
}

func NewEntryIndex(index *dindex.DomeIndex, pathCtx *pathctx.PathCtx) *EntryIndex {
	return &EntryIndex{
		index:   index,
		pathCtx: pathCtx,
		log:     logrus.WithField("service", "entry_index"),
	}
}

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 4096

func (mw *EntryIndex) GetMimeType(realpath string) (entrymodel.MimeType, error) {
	ext := filepath.Ext(realpath)
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		return entrymodel.MimeType(ctype), nil
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

	return entrymodel.MimeType(ctype), nil
}

func (mw *EntryIndex) GetEntryInfo(ctx context.Context, virtpath string) (*entrymodel.EntryInfo, error) {
	return mw.genEntryInfo(ctx, mw.pathCtx.MotherPath(virtpath), virtpath)
}

func (mw *EntryIndex) GetPreviewsDir(ctx context.Context, virtpath string) (domefile.File, error) {
	realpath := mw.pathCtx.MotherPath(virtpath)
	info, err := mw.genEntryInfo(ctx, realpath, virtpath)
	if err != nil {
		mw.log.Errorf("Error generating video info: %w", err)
		return nil, err
	}

	err = mw.genPreviews(ctx, realpath, virtpath, getTimestamps(info.MediaInfo.VideoInfo.Duration))
	if err != nil {
		mw.log.Errorf("Eror generating video previews: %w", err)
		return nil, err
	}

	file, err := os.Open(mw.getPreviewsDirPath(virtpath))
	if err != nil {
		return nil, err
	}
	return domefile.WrapOsFile("previews/", file), nil
}

func (mw *EntryIndex) GetPreviewFile(ctx context.Context, virtpath string, previewFileName string) (domefile.File, error) {
	mediapath := mw.pathCtx.MotherPath(virtpath)

	err := mw.generateCache(ctx, mediapath, virtpath)
	if err != nil {
		return nil, err
	}

	medianame := filepath.Base(mediapath)
	file, err := os.Open(path.Join(mw.getPreviewsDirPath(virtpath), previewFileName))
	if err != nil {
		return nil, err
	}
	return domefile.WrapOsFile(medianame, file), nil
}
