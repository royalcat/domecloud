package domefs

import (
	"context"
	"dmch-server/src/domefs/domefile"
	"encoding/json"
	"os"
	"path"
)

type vServeFile func(name string) (domefile.File, error) // returns "", nil  if must be sorted as normal file

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

func (domefs *DomeFS) wrapOpenFile(virtpath string) (domefile.File, error) {
	osfile, err := os.Open(domefs.pathctx.MotherPath(virtpath))
	if err != nil {
		return nil, err
	}
	return domefile.WrapOsFile(path.Base(virtpath), osfile), nil
}

func (domefs *DomeFS) serveVirtualEntryToReal(virtpath string) (domefile.File, error) {
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

	return domefs.wrapOpenFile(virtpath)
}

func (domefs *DomeFS) serveInfoJson(virtPath string) (domefile.File, error) {
	ctx := context.Background()
	virtFilePath := path.Dir(virtPath)
	entryInfo, err := domefs.cache.GetEntryInfo(ctx, virtFilePath)
	if err != nil {
		return nil, err
	}

	entryInfoJson, err := json.Marshal(entryInfo)
	if err != nil {
		return nil, err
	}
	return domefile.NewMemoryFile("info.json", "application/json", entryInfoJson), nil
}

func (domefs *DomeFS) servePreviewsDir(virtPath string) (domefile.File, error) {
	ctx := context.Background()
	virtFilePath := path.Dir(virtPath)
	dirFile, err := domefs.cache.GetPreviewsDir(ctx, virtFilePath)
	if err != nil {
		return nil, err
	}
	return dirFile, nil
}

func (domefs *DomeFS) servePreview(virtPath string) (domefile.File, error) {
	ctx := context.Background()
	virtFilePath := path.Dir(path.Dir(virtPath))
	dirFile, err := domefs.cache.GetPreviewsDir(ctx, virtFilePath)
	if err != nil {
		return nil, err
	}
	return dirFile, nil
}
