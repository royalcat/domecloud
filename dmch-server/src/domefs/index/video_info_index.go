package indexes

import (
	"dmch-server/src/domefs/media"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/bradenaw/juniper/container/tree"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

type VideoInfoIndex struct {
	mutex sync.RWMutex

	pathToInfo map[string]media.VideoInfo

	durationIndex tree.Map[time.Duration, string]
	//pathTreeIndex
}

func NewVideoInfoIndex() *VideoInfoIndex {
	return &VideoInfoIndex{
		durationIndex: tree.NewMap[time.Duration, string](xsort.OrderedLess[time.Duration]),
	}
}

func (vii *VideoInfoIndex) Set(v media.VideoInfo) {
	vii.mutex.Lock()

	vii.pathToInfo[v.Path] = v
	vii.durationIndex.Put(v.Duration, path.Clean(v.Path))

	vii.mutex.Unlock()
}

func (vii *VideoInfoIndex) GetSortedByDuration(targetDir string, recursive bool) []media.VideoInfo {
	targetDir = path.Clean(targetDir)
	var filter func(tree.KVPair[time.Duration, string]) bool
	if recursive {
		filter = func(t tree.KVPair[time.Duration, string]) bool {
			return strings.HasPrefix(t.Value, targetDir)
		}
	} else {
		filter = func(t tree.KVPair[time.Duration, string]) bool {
			match, err := path.Match(targetDir+"/*", t.Value)
			return err != nil && match
		}
	}

	vii.mutex.Lock()
	defer vii.mutex.Unlock()

	return iterator.Collect(
		iterator.Map(
			iterator.Filter(
				vii.durationIndex.Iterate(),
				filter,
			),
			func(t tree.KVPair[time.Duration, string]) media.VideoInfo {
				return vii.pathToInfo[t.Value]
			},
		),
	)
}
