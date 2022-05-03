package domefs

import (
	"dmch-server/src/domefs/media"
	"dmch-server/src/domefs/mediacache"
	"io"
	"io/fs"
	"os"
	"path"

	"go.mongodb.org/mongo-driver/mongo"
)

type DomeFS struct {
	rootDir string
	//cacheDir string

	vfuncVirtFile []map[string]vServeFile // [level][name]Serve

	cache *mediacache.MediaCache
}

func NewDomeFS(db *mongo.Database, rootDir, cacheDir string) *DomeFS {
	dfs := &DomeFS{
		rootDir: path.Clean(rootDir),

		cache: mediacache.NewMediaCache(db, path.Clean(cacheDir)),
	}
	dfs.initVirtFileFunctions()
	return dfs
}

func (domefs *DomeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	realpath, err := domefs.RealPath(name)
	if err != nil {
		return nil, err
	}

	return os.ReadDir(realpath)
}

func (domefs *DomeFS) Open(name string) (File, error) {
	realpath, err := domefs.RealPath(name)
	if err != nil {
		return nil, err
	}
	return os.Open(realpath)
}

func (domefs *DomeFS) Stat(name string) (fs.FileInfo, error) {
	realpath, err := domefs.RealPath(name)
	if err != nil {
		return nil, err
	}
	return os.Stat(realpath)
}

func (domefs DomeFS) RealPath(fpath string) (string, error) {
	fpath, err := domefs.serveVirtualEntryToReal(fpath)
	if err != nil {
		return "", err
	}
	return fpath, nil
}

func (domefs *DomeFS) MimeType(name string) (media.MimeType, error) {
	realpath, err := domefs.RealPath(name)
	if err != nil {
		return "", err
	}
	return domefs.cache.GetMimeType(realpath)
}

type File interface {
	fs.File
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
}
