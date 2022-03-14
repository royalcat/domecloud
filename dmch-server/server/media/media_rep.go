package media

import (
	"context"
	"dmch-server/cfs"
	"errors"
	"fmt"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
)

type MediaInfoRepository struct {
	videoCacheMutex sync.Mutex
	video           map[string]*VideoInfo
}

func NewMediaCache() (*MediaInfoRepository, error) {
	return &MediaInfoRepository{
		video: make(map[string]*VideoInfo),
	}, nil
}

var ErrVideoStreamNotFound = errors.New("Video stream not found")

func (w *MediaInfoRepository) GenerateCache(ctx context.Context, dirpath string) {
	entries, err := cfs.Cfs.ReadDir(dirpath)
	if err != nil {
		logrus.Errorf("Failed to read dir %s with error: %s", dirpath, err.Error())
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			fpath := path.Join(dirpath, entry.Name())
			w.GetVideo(ctx, fpath)
		}
	}
}

func (rep *MediaInfoRepository) ListVideo(ctx context.Context, dirpath string) ([]*VideoInfo, error) {
	entries, err := cfs.Cfs.ReadDir(dirpath)
	if err != nil {
		return []*VideoInfo{}, fmt.Errorf("Failed to read dir %s with error: %s", dirpath, err.Error())
	}
	videos := []*VideoInfo{}
	for _, entry := range entries {
		if !entry.IsDir() {
			fpath := path.Join(dirpath, entry.Name())
			if video, err := rep.GetVideo(ctx, fpath); err == nil {
				videos = append(videos, video)
			}

		}
	}
	return videos, nil
}

func (rep *MediaInfoRepository) GetVideo(ctx context.Context, fpath string) (*VideoInfo, error) {
	rep.videoCacheMutex.Lock()
	info, ok := rep.video[fpath]
	if !ok {
		info, err := GenerateVideoInfo(ctx, fpath)
		if err != nil {
			logrus.Errorf("Failed generate video info: %s with error: %s", fpath, err.Error())
			return nil, err
		}
		rep.video[fpath] = info
	}
	rep.videoCacheMutex.Unlock()
	return info, nil
}
