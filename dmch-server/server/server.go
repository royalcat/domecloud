package server

import (
	"context"
	"dmch-server/server/media"
)

type DmServer struct {
	mediaRepository *media.MediaInfoRepository
}

func (d *DmServer) ListVideos(ctx context.Context, dirpath string) ([]*media.VideoInfo, error) {
	return d.mediaRepository.ListVideo(ctx, dirpath)
}
