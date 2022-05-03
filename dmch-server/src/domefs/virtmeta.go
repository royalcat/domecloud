package domefs

import (
	"context"
	"dmch-server/src/domefs/media"
	"errors"
	"path"
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

func (domefs *DomeFS) serveVirtualEntryToReal(virtpath string) (string, error) {
	namePart := virtpath
	for _, funcMap := range domefs.vfuncVirtFile {
		cmd := path.Base(namePart)
		for key, fun := range funcMap {
			if cmd == key {
				return fun(virtpath)
			}
		}
		namePart = path.Dir(namePart)
	}

	return domefs.rootJoinedPath(virtpath), nil
}

func (domefs *DomeFS) serveInfoJson(originalName string) (string, error) {
	ctx := context.Background()
	virtFilePath := path.Dir(originalName)

	mimetype, err := domefs.MimeType(path.Base(virtFilePath))
	if err != nil {
		return "", err
	}
	if mimetype.MediaType() != media.MediaTypeVideo {
		return "", errors.New("invalid file media type")
	}

	return domefs.cache.GetInfoFilePath(ctx, domefs.rootJoinedPath(virtFilePath), virtFilePath)
}

func (domefs *DomeFS) servePreviewsDir(originalName string) (string, error) {
	ctx := context.Background()
	virtFilePath := path.Dir(originalName)
	return domefs.cache.GetPreviewsDirPath(ctx, domefs.rootJoinedPath(virtFilePath), virtFilePath)
}

func (domefs *DomeFS) servePreview(originalName string) (string, error) {
	ctx := context.Background()
	virtFilePath := path.Dir(path.Dir(originalName))
	previewCachePath, err := domefs.cache.GetPreviewsDirPath(ctx, domefs.rootJoinedPath(virtFilePath), virtFilePath)
	if err != nil {
		return "", err
	}
	return path.Join(previewCachePath, path.Base(originalName)), nil
}

func (domefs *DomeFS) rootJoinedPath(virtpath string) string {
	// if virtpath == "/" {
	// 	return domefs.rootDir
	// }
	return path.Join(domefs.rootDir, virtpath)
}
