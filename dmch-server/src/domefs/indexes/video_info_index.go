package indexes

import (
	"dmch-server/src/domefs/media"
	"strings"
	"sync"
	"time"

	"github.com/puzpuzpuz/xsync"
)

type VideoInfoIndex struct {
	mutex sync.RWMutex

	maps *xsync.MapOf[media.VideoInfo]

	pathToInfo map[string]media.VideoInfo

	durationIndex SortedMap[time.Duration, string]
}

func NewVideoInfoIndex() *VideoInfoIndex {
	infoindex := &VideoInfoIndex{
		maps: xsync.NewMapOf[media.VideoInfo](),
	}
	return infoindex
}

func (vii *VideoInfoIndex) Set(v media.VideoInfo) {
	vii.mutex.Lock()

	vii.pathToInfo[v.Path] = v
	vii.durationIndex.Set(v.Duration, v.Path)

	vii.mutex.Unlock()
}

func (vii *VideoInfoIndex) GetSortedByDuration(dir string) []media.VideoInfo {
	out := make([]media.VideoInfo, 0, len(vii.pathToInfo)/2)

	for _, path := range vii.durationIndex.All() {
		if strings.HasPrefix(path, dir) {
			out = append(out, vii.pathToInfo[path])
		}
	}

	return out
}
