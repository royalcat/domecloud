package indexapi

import (
	"context"
	"dmch-server/src/domefs/dindex"
	"dmch-server/src/domefs/entrymodel"
)

type DomeIndexApi struct {
	index *dindex.DomeIndex
}

func NewIndexApi(index *dindex.DomeIndex) *DomeIndexApi {
	return &DomeIndexApi{
		index: index,
	}
}

func (dia *DomeIndexApi) ListMedia(ctx context.Context, targetDir string) ([]entrymodel.EntryInfo, error) {
	return dia.index.VideoInfo.GetMediaInDir(ctx, targetDir, true)
}
