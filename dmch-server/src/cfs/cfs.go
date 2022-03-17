package cfs

import (
	"context"
	"dmch-server/src/config"

	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

type DmFS struct {
	rootDir  string
	cacheDir string
}

func NewDmFS() *DmFS {
	return &DmFS{
		rootDir:  path.Clean(config.Config.RootFolder),
		cacheDir: path.Clean(config.Config.CacheFolder),
	}
}

func (dmfs *DmFS) ReadDir(name string) ([]fs.DirEntry, error) {
	cmd := path.Base(name)
	if cmd == "previews" || cmd == "info" { // TODO седелать мапу с разными функциями
		dir := path.Dir(name)
		if EndsWithOneOf(dir, config.Config.Media.Extensions) { // BUG не брать тут из конфига
			switch cmd {
			case "previews":
				ctx := context.Background()
				info, _ := dmfs.getVideoInfo(ctx, name)
				dmfs.getPreviews(ctx, name, getTimestamps(info.Duration))
				return os.ReadDir(dmfs.getPreviewsRealPath(name))
			}
		}
	}

	return os.ReadDir(dmfs.RealPath(name))
}

// Open implements fs.StatFS
func (dmfs *DmFS) Open(name string) (fs.File, error) {
	cmd := path.Base(name)
	if cmd == "previews" || cmd == "info" { // TODO седелать мапу с разными функциями
		fpath := path.Dir(name)
		if EndsWithOneOf(strings.ToLower(fpath), config.Config.Media.Extensions) { // BUG не брать тут из конфига
			switch cmd {
			case "previews":
				ctx := context.Background()
				info, _ := dmfs.getVideoInfo(ctx, fpath)
				dmfs.getPreviews(ctx, fpath, getTimestamps(info.Duration))
				return os.Open(dmfs.getPreviewsRealPath(fpath))
			case "info":
				ctx := context.Background()
				_, err := dmfs.getVideoInfo(ctx, fpath)
				if err != nil {
					logrus.Errorf("error getting video info: %s", err.Error())
				}
				infoPath := dmfs.getInfoRealPath(fpath)
				return os.Open(infoPath)
			}
		}
	} else if path.Base(path.Dir(name)) == "previews" {
		fpath := path.Dir(path.Dir(name))
		ctx := context.Background()
		info, _ := dmfs.getVideoInfo(ctx, fpath)
		dmfs.getPreviews(ctx, fpath, getTimestamps(info.Duration))
		return os.Open(path.Join(dmfs.getPreviewsRealPath(fpath), path.Base(name)))
	}

	return os.Open(dmfs.RealPath(name))
}

// Stat implements fs.StatFS
func (dmfs *DmFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (dmfs DmFS) RealPath(fpath string) string {
	return path.Join(dmfs.rootDir, fpath)
}

var _ fs.StatFS = (*DmFS)(nil)
var _ fs.ReadDirFS = (*DmFS)(nil)

func EndsWithOneOf(s string, ends []string) bool {
	for _, end := range ends {
		if strings.HasSuffix(s, end) {
			return true
		}
	}
	return false
}

// type PreviewVFile struct {
// 	file     *os.File
// 	realpath string
// }

// // Info implements fs.DirEntry
// func (f *PreviewVFile) Info() (fs.FileInfo, error) {
// 	fs.Dir
// }

// // IsDir implements fs.DirEntry
// func (*PreviewVFile) IsDir() bool {
// 	return false
// }

// // Name implements fs.DirEntry
// func (f *PreviewVFile) Name() string {
// 	return path.Base(f.realpath)
// }

// // Type implements fs.DirEntry
// func (*PreviewVFile) Type() fs.FileMode {
// 	panic("unimplemented")
// }

// var _ fs.DirEntry = (*PreviewVFile)(nil)
