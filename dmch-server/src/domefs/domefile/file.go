package domefile

import (
	"dmch-server/src/domefs/dmime"
	"io"
	"io/fs"
	"mime"
	"os"
	"path"
)

type File interface {
	Name() string
	Path() string
	fs.DirEntry
	fs.File
	io.Seeker
	ReadDir(n int) ([]DirEntry, error)
	//Stat() (fs.FileInfo, error)
	MimeType
}

type DirEntry interface {
	fs.DirEntry
	MimeType
}

type MimeType interface {
	MimeType() dmime.MimeType
}

type wrapFile struct {
	osfile   *os.File
	mimeType dmime.MimeType
}

func OpenDomeFile(fpath string) (File, error) {
	fi, err := os.Lstat(fpath)
	if err != nil {
		return nil, err
	}
	realpath := fpath
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		realpath, err = os.Readlink(fpath)
		if err != nil {
			return nil, err
		}
	}

	osfile, err := os.Open(realpath)
	if err != nil {
		return nil, err
	}
	return WrapOsFile(path.Base(fpath), osfile), nil
}

func WrapOsFile(name string, f *os.File) File {

	var mimeType dmime.MimeType
	ext := path.Ext(name)
	if name[len(name)-1] == '/' || ext == "" {
		mimeType = dmime.MimeTypeDirectory
	} else {
		mimeType = dmime.MimeType(mime.TypeByExtension(ext))
	}

	return &wrapFile{
		osfile:   f,
		mimeType: mimeType,
	}
}

// Path implements File
func (f *wrapFile) Path() string {
	return f.osfile.Name()
}

// Info implements File
func (f *wrapFile) Info() (fs.FileInfo, error) {
	return f.osfile.Stat()
}

// IsDir implements File
func (f *wrapFile) IsDir() bool {
	return f.mimeType == dmime.MimeTypeDirectory
}

// Type implements File
func (f *wrapFile) Type() fs.FileMode {
	return f.Type()
}

// Name implements File
func (f *wrapFile) Name() string {
	return path.Base(f.osfile.Name())
}

// Close implements File
func (f *wrapFile) Close() error {
	return f.osfile.Close()
}

// Read implements File
func (f *wrapFile) Read(b []byte) (int, error) {
	return f.osfile.Read(b)
}

// Stat implements File
func (f *wrapFile) Stat() (fs.FileInfo, error) {
	return f.osfile.Stat()
}

// Seek implements File
func (f *wrapFile) Seek(offset int64, whence int) (int64, error) {
	return f.osfile.Seek(offset, whence)
}

// MimeType implements File
func (f *wrapFile) MimeType() dmime.MimeType {
	return f.mimeType
}

func (f *wrapFile) ReadDir(n int) ([]DirEntry, error) {
	entries, err := f.osfile.ReadDir(n)
	if err != nil {
		return nil, err
	}

	domeEntries := make([]DirEntry, 0, len(entries))
	for _, entry := range entries {
		domeEntries = append(domeEntries, WrapDomeEntry(f.osfile.Name(), entry))
	}

	return domeEntries, nil
}

type domeEntry struct {
	fs.DirEntry
}

func WrapDomeEntry(dir string, e fs.DirEntry) DirEntry {
	if e.Type() == fs.ModeSymlink {
		entry, _ := OpenDomeFile(path.Join(dir, e.Name()))
		return entry
	}

	return domeEntry{e}
}

// MimeType implements DirEntry
func (e domeEntry) MimeType() dmime.MimeType {

	if e.IsDir() {
		return dmime.MimeTypeDirectory
	}

	return dmime.MimeType(mime.TypeByExtension(path.Ext(e.Name())))
}
