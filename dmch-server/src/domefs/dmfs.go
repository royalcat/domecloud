package domefs

import (
	"dmch-server/src/config"
	"io"
	"io/fs"
	"os"
	"path"
)

type DomeFS struct {
	rootDir  string
	cacheDir string

	vfuncVirtFile []map[string]vServeFile // [level][name]Serve
}

func NewDomeFS() *DomeFS {
	dfs := &DomeFS{
		rootDir:  path.Clean(config.Config.RootFolder),
		cacheDir: path.Clean(config.Config.CacheFolder),
	}
	dfs.initVirtFileFunctions()
	return dfs
}

func (domefs *DomeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	// cmd := path.Base(name)
	// if cmd == "previews" || cmd == "info" { // TODO седелать мапу с разными функциями
	// 	fpath := path.Dir(name)
	// 	mimetype, _ := domefs.MimeType(path.Base(fpath))
	// 	if mimetype.MediaType() == media.MediaTypeVideo {
	// 		switch cmd {
	// 		case "previews":
	// 			ctx := context.Background()
	// 			info, _ := domefs.getVideoInfo(ctx, name)
	// 			domefs.getPreviews(ctx, name, getTimestamps(info.Duration))
	// 			return os.ReadDir(domefs.getPreviewsRealPath(name))
	// 		}
	// 	}
	// }

	fpath, err := domefs.serveVirtualEntry(name)
	if fpath != "" && err == nil {
		return os.ReadDir(fpath)
	} else if err != nil {
		return nil, err
	}

	return os.ReadDir(domefs.realPath(name))
}

// Open implements fs.StatFS
func (domefs *DomeFS) Open(name string) (File, error) {
	// cmd := path.Base(name)
	// if cmd == "previews" || cmd == "info" {
	// 	fpath := path.Dir(name)
	// 	if stat, err := domefs.Stat(fpath); err == nil && stat.IsDir() {
	// 		if stat, err = domefs.Stat(name); err == nil { // return normal file if exists
	// 			return os.Open(domefs.RealPath(name))
	// 		} else {
	// 			return nil, fs.ErrInvalid
	// 		}
	// 	}

	// 	mimetype, _ := domefs.MimeType(path.Base(fpath))
	// 	if mimetype.MediaType() == media.MediaTypeVideo {
	// 		switch cmd {
	// 		case "previews":
	// 			ctx := context.Background()
	// 			info, _ := domefs.getVideoInfo(ctx, fpath)
	// 			domefs.getPreviews(ctx, fpath, getTimestamps(info.Duration))
	// 			return os.Open(domefs.getPreviewsRealPath(fpath))
	// 		case "info":
	// 			ctx := context.Background()
	// 			_, err := domefs.getVideoInfo(ctx, fpath)
	// 			if err != nil {
	// 				logrus.Errorf("error getting video info: %s", err.Error())
	// 			}
	// 			infoPath := domefs.getInfoRealPath(fpath)
	// 			return os.Open(infoPath)
	// 		}
	// 	} else {
	// 		return nil, fs.ErrInvalid
	// 	}
	// } else if path.Base(path.Dir(name)) == "previews" {
	// 	fpath := path.Dir(path.Dir(name))
	// 	ctx := context.Background()
	// 	info, _ := domefs.getVideoInfo(ctx, fpath)
	// 	domefs.getPreviews(ctx, fpath, getTimestamps(info.Duration))
	// 	return os.Open(path.Join(domefs.getPreviewsRealPath(fpath), path.Base(name)))
	// }

	fpath, err := domefs.serveVirtualEntry(name)
	if fpath != "" && err == nil {
		return os.Open(fpath)
	} else if err != nil {
		return nil, err
	}

	return os.Open(domefs.realPath(name))
}

func (domefs *DomeFS) Stat(name string) (fs.FileInfo, error) {
	fpath, err := domefs.serveVirtualEntry(name)
	if fpath != "" && err == nil {
		return os.Stat(fpath)
	} else if err != nil {
		return nil, err
	}

	realPath := domefs.realPath(name)
	return os.Stat(realPath)
}

func (domefs DomeFS) realPath(fpath string) string {
	return path.Join(domefs.rootDir, fpath)
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
