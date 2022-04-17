package indexes

import (
	"dmch-server/src/domefs/media"
	"path"
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
	vii.durationIndex.Set(v.Duration, path.Clean(v.Path))

	vii.mutex.Unlock()
}

func (vii *VideoInfoIndex) GetSortedByDuration(targetDir string, recursive bool) []media.VideoInfo {
	out := make([]media.VideoInfo, 0, len(vii.pathToInfo)/2)
	targetDir = path.Clean(targetDir)

	if recursive {
		for _, videopath := range vii.durationIndex.All() {
			if strings.HasPrefix(videopath, targetDir) {
				out = append(out, vii.pathToInfo[videopath])
			}
		}
	} else {
		for _, videopath := range vii.durationIndex.All() {
			if match, err := path.Match(targetDir+"/*", videopath); err != nil && match {
				out = append(out, vii.pathToInfo[videopath])
			}
		}
	}

	return out
}
