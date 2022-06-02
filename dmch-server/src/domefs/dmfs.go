package domefs

import (
	"dmch-server/src/domefs/dindex"
	"dmch-server/src/domefs/domefile"
	"dmch-server/src/domefs/entrycache"
	"dmch-server/src/domefs/indexapi"
	"dmch-server/src/domefs/pathctx"
	"io/fs"

	"go.mongodb.org/mongo-driver/mongo"
)

type DomeFS struct {
	pathctx *pathctx.PathCtx

	index *dindex.DomeIndex

	vfuncVirtFile []map[string]vServeFile // [level][name]Serve

	cache *entrycache.EntryIndex

	Api *indexapi.DomeIndexApi
}

func NewDomeFS(db *mongo.Database, rootDir, cacheDir string) *DomeFS {
	index := dindex.NewDomeIndex(db)
	pathctx := pathctx.NewPathCtx(rootDir, cacheDir)
	dfs := &DomeFS{
		pathctx: pathctx,

		cache: entrycache.NewEntryIndex(index, pathctx),

		Api: indexapi.NewIndexApi(index),
	}
	dfs.initVirtFileFunctions()
	return dfs
}

func (domefs *DomeFS) ReadDir(name string) ([]domefile.DirEntry, error) {
	file, err := domefs.serveVirtualEntryToReal(name)
	if err != nil {
		return nil, err
	}

	return file.ReadDir(0)
}

func (domefs *DomeFS) Open(name string) (domefile.File, error) {
	return domefs.serveVirtualEntryToReal(name)
}

func (domefs *DomeFS) Stat(name string) (fs.FileInfo, error) {
	realpath, err := domefs.serveVirtualEntryToReal(name)
	if err != nil {
		return nil, err
	}
	return realpath.Stat()
}
