package cfs

import (
	"dmch-server/config"
	"io/fs"
	"path"
)

var Cfs *DmFS = &DmFS{RootPath: config.Config.RootFolder}

type DmFS struct {
	RootPath string
}

// ReadDir implements fs.ReadDirFS
func (*DmFS) ReadDir(name string) ([]fs.DirEntry, error) {

	panic("unimplemented")
}

func NewCfs(rootPath string) {}

func (dmfs DmFS) RealPath(fpath string) string {
	return path.Join(dmfs.RootPath, fpath)
}

// Open implements fs.StatFS
func (*DmFS) Open(name string) (fs.File, error) {
	panic("unimplemented")
}

// Stat implements fs.StatFS
func (*DmFS) Stat(name string) (fs.FileInfo, error) {
	panic("unimplemented")
}

var _ fs.StatFS = (*DmFS)(nil)
var _ fs.ReadDirFS = (*DmFS)(nil)
