package mediacache

import (
	"context"
	"dmch-server/src/domefs"
	"dmch-server/src/domefs/media"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
)

type MediaCache struct {
	dfs *domefs.DomeFS

	fflock sync.Mutex

	cacheDir string

	log *logrus.Entry
}

func (mw *MediaCache) GetInfoFilePath(ctx context.Context, virtpath string) (string, error) {
	realpath := mw.dfs.RealPath(virtpath)
	_, err := mw.genVideoInfo(ctx, virtpath, realpath)
	if err != nil {
		mw.log.Errorf("Eror generating video info: %w", err)
		return "", err
	}

	return mw.getInfoPath(virtpath), nil
}

func (mw *MediaCache) GetPreviewsDirPath(ctx context.Context, virtpath string) (string, error) {
	realpath := mw.dfs.RealPath(virtpath)
	info, err := mw.genVideoInfo(ctx, virtpath, realpath)
	if err != nil {
		mw.log.Errorf("Eror generating video info: %w", err)
		return "", err
	}

	err = mw.genPreviews(ctx, virtpath, realpath, getTimestamps(info.Duration))
	if err != nil {
		mw.log.Errorf("Eror generating video previews: %w", err)
		return "", err
	}

	return mw.getPreviewsDirPath(virtpath), nil
}

func (mw *MediaCache) RunWarmer() {

	//go mw.worker()
}

func (mw *MediaCache) warmDirectory(dir string, recursive bool) error {
	tasks, err := mw.listDirMedia(dir, recursive)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		mw.generateCache(context.Background(), task)
	}

	return nil
}

func (mw *MediaCache) listDirMedia(dir string, recursive bool) ([]string, error) {
	medias := []string{}

	entries, err := mw.dfs.ReadDir(dir)
	if err != nil {
		return medias, err
	}
	for _, entry := range entries {
		entrypath := path.Join(dir, entry.Name())
		if entry.IsDir() {
			list, err := mw.listDirMedia(entrypath, recursive)
			if err != nil {
				return medias, err
			}
			medias = append(medias, list...)
		} else if mt, err := mw.dfs.MimeType(entrypath); err != nil && mt.MediaType() != media.MediaTypeNone {
			medias = append(medias, entrypath)
		}
	}

	return medias, nil
}
