package domefs

import (
	"context"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type VServe func(name string) (File, error)

func (domefs *DomeFS) initVirtFileFunctions() {}

func (domefs *DomeFS) serveVirtualEntry(name string) (File, error) {

	namePart := name
	for _, funcMap := range domefs.vfuncVirtFile {
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

func (domefs *DomeFS) serveInfoJson(originalName string) (File, error) {
	ctx := context.Background()
	fpath := path.Dir(originalName)
	_, err := domefs.getVideoInfo(ctx, fpath)
	if err != nil {
		logrus.Errorf("error getting video info: %s", err.Error())
	}
	infoPath := domefs.getInfoRealPath(fpath)
	return os.Open(infoPath)
}
