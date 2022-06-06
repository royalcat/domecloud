package pathctx

import "path"

type PathCtx struct {
	rootDir   string
	cachePath string
}

func NewPathCtx(rootDir, cachePath string) *PathCtx {
	return &PathCtx{rootDir: rootDir, cachePath: cachePath}
}

func (pathctx *PathCtx) MotherPath(virtpath string) string {
	if virtpath == "/" {
		return pathctx.rootDir
	}
	return path.Join(pathctx.rootDir, virtpath)
}

func (pathctx *PathCtx) CachePath(virtpath string) string {
	if virtpath == "/" {
		return pathctx.rootDir
	}
	return path.Join(pathctx.cachePath, virtpath)
}
