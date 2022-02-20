package cfs

import (
	"errors"
	"io/fs"
	"os"
	"runtime"
)

type Cfs struct {
	dir string
}

var _ fs.FS = (*Cfs)(nil)
var _ fs.ReadDirFS = (*Cfs)(nil)
var _ fs.StatFS = (*Cfs)(nil)

func NewCfs(dir string) (*Cfs, error) {
	if runtime.GOOS == "windows" {
		return nil, errors.New("windows not supported")
	}

	cfs := &Cfs{
		dir: dir,
	}

	return cfs, nil
}

func (cfs *Cfs) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	return os.Open(string(cfs.dir) + "/" + name)
}

func (cfs *Cfs) ReadDir(name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}

	return os.ReadDir(string(cfs.dir) + "/" + name)
}

func (cfs *Cfs) Stat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrInvalid}
	}
	return os.Stat(string(cfs.dir) + "/" + name)
}

func (cfs *Cfs) Mkdir(name string, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrInvalid}
	}
	return os.Mkdir(name, perm)
}
