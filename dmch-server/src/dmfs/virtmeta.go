package dmfs

import (
	"context"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type VServe func(name string) (File, error)

func (dmfs *DmFS) initVirtFileFunctions() {}

func (dmfs *DmFS) serveVirtualEntry(name string) (File, error) {

	namePart := name
	for _, funcMap := range dmfs.vfuncVirtFile {
		cmd := path.Base(namePart)
		for key, fun := range funcMap {
			if cmd == key {
				return fun(name)
			}
		}
		namePart = path.Dir(namePart)
	}

	return nil, nil
}

func (dmfs *DmFS) serveInfoJson(originalName string) (File, error) {
	ctx := context.Background()
	fpath := path.Dir(originalName)
	_, err := dmfs.getVideoInfo(ctx, fpath)
	if err != nil {
		logrus.Errorf("error getting video info: %s", err.Error())
	}
	infoPath := dmfs.getInfoRealPath(fpath)
	return os.Open(infoPath)
}
