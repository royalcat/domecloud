package domefs

import (
	"dmch-server/src/domefs/media"
	"errors"
	"path"

	"github.com/sirupsen/logrus"
)

type vServeFile func(name string) (string, error) // returns "", nil  if must be sorted as normal file

func (domefs *DomeFS) initVirtFileFunctions() {
	domefs.vfuncVirtFile = []map[string]vServeFile{
		{
			"info.json": domefs.serveInfoJson,
			"previews":  domefs.servePreviewsDir,
		},
		{
			"previews": domefs.servePreview,
		},
	}
}

func (domefs *DomeFS) serveVirtualEntry(name string) (string, error) {
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

	return "", nil
}

func (domefs *DomeFS) serveInfoJson(originalName string) (string, error) {
	fpath := path.Dir(originalName)

	mimetype, err := domefs.MimeType(path.Base(fpath))
	if err != nil {
		return "", err
	}
	if mimetype.MediaType() != media.MediaTypeVideo {
		return "", errors.New("invalid file media type")
	}

	_, err = domefs.getVideoInfo(fpath)
	if err != nil {
		logrus.Errorf("error getting video info: %s", err.Error())
		return "", err
	}
	infoPath := domefs.getInfoRealPath(fpath)

	return infoPath, nil
}

func (domefs *DomeFS) servePreviewsDir(originalName string) (string, error) {
	fpath := path.Dir(originalName)
	info, err := domefs.getVideoInfo(fpath)
	if err != nil {
		logrus.Errorf("error getting video info: %s", err.Error())
		return "", err
	}
	domefs.getPreviews(fpath, getTimestamps(info.Duration))
	return domefs.getPreviewsRealPath(fpath), nil
}

func (domefs *DomeFS) servePreview(originalName string) (string, error) {
	fpath := path.Dir(path.Dir(originalName))
	info, err := domefs.getVideoInfo(fpath)
	if err != nil {
		logrus.Errorf("error getting video info: %s", err.Error())
		return "", err
	}
	domefs.getPreviews(fpath, getTimestamps(info.Duration))
	return path.Join(domefs.getPreviewsRealPath(fpath), path.Base(originalName)), nil
}
