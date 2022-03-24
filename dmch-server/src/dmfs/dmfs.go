package dmfs

import (
	"context"
	"dmch-server/src/config"
	"dmch-server/src/dmfs/media"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type DmFS struct {
	rootDir  string
	cacheDir string

	vfuncVirtFile []map[string]VServe // [level][name]Serve
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
		fpath := path.Dir(name)
		mimetype, _ := dmfs.MimeType(path.Base(fpath))
		if mimetype.MediaType() == media.MediaTypeVideo { // BUG не брать тут из конфига
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
func (dmfs *DmFS) Open(name string) (File, error) {

	cmd := path.Base(name)
	if cmd == "previews" || cmd == "info" { // TODO седелать мапу с разными функциями
		fpath := path.Dir(name)
		if stat, err := dmfs.Stat(fpath); err == nil && stat.IsDir() {
			if stat, err = dmfs.Stat(name); err == nil { // return normal file if exists
				return os.Open(dmfs.RealPath(name))
			} else {
				return nil, fs.ErrInvalid
			}
		}

		mimetype, _ := dmfs.MimeType(path.Base(fpath))
		if mimetype.MediaType() == media.MediaTypeVideo { // BUG не брать тут из конфига
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
		} else {
			return nil, fs.ErrInvalid
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
	return os.Stat(dmfs.RealPath(name))
}

func (dmfs DmFS) RealPath(fpath string) string {
	return path.Join(dmfs.rootDir, fpath)
}

// var _ fs.StatFS = (*DmFS)(nil)
// var _ fs.ReadDirFS = (*DmFS)(nil)

type File interface {
	fs.File
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
}
