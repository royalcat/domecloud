package mediacache

import (
	"os"
	"path"
)

func (mw *MediaCache) RunWarmer() {

}

func (mw *MediaCache) warmDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			mw.warmDirectory(path.Join(dir, entry.Name()))
		}

	}

	return nil
}

func (mw *MediaCache) warm(file string) {

}

type genCacheTask struct {
	virtpath, realpath string
}
