package mediacache

import "path"

func (mw *MediaCache) getPreviewsDirPath(virtpath string) string {
	return path.Join(mw.cacheDir, virtpath, "previews")
}

func (mw *MediaCache) getInfoPath(virtpath string) string {
	return path.Join(mw.cacheDir, virtpath, "info.json")
}
