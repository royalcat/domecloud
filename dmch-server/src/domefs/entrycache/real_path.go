package entrycache

import "path"

func (mw *EntryIndex) getPreviewsDirPath(virtpath string) string {
	return mw.pathCtx.CachePath(path.Join(virtpath, "previews"))
}

func (mw *EntryIndex) getInfoPath(virtpath string) string {
	return mw.pathCtx.CachePath(path.Join(virtpath, "info.json"))
}
