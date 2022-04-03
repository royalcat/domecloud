package mediacache

import (
	"context"

	"github.com/sirupsen/logrus"
)

type MediaCache struct {
	cacheDir string

	log *logrus.Entry
}

func (mw *MediaCache) GetInfoFilePath(ctx context.Context, virtpath, realpath string) (string, error) {
	_, err := mw.genVideoInfo(ctx, virtpath, realpath)
	if err != nil {
		mw.log.Errorf("Eror generating video info: %w", err)
		return "", err
	}

	return mw.getInfoPath(virtpath), nil
}

func (mw *MediaCache) GetPreviewsDirPath(ctx context.Context, virtpath, realpath string) (string, error) {
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
