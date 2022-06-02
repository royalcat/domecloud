package domefile

import (
	"dmch-server/src/domefs/entrymodel"
	"io"
	"io/fs"
	"mime"
	"os"
	"path"
)

type File interface {
	Name() string
	fs.File
	io.Seeker
	ReadDir(n int) ([]DirEntry, error)
	Stat() (fs.FileInfo, error)
	MimeType
}

type DirEntry interface {
	fs.DirEntry
	MimeType
}

type MimeType interface {
	MimeType() entrymodel.MimeType
}

type wrapFile struct {
	osfile   *os.File
	mimeType entrymodel.MimeType
}

func WrapOsFile(name string, f *os.File) File {
	return &wrapFile{
		osfile:   f,
		mimeType: entrymodel.MimeType(mime.TypeByExtension(path.Ext(name))),
	}
}

// Name implements File
func (f *wrapFile) Name() string {
	return f.osfile.Name()
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
func (f *wrapFile) MimeType() entrymodel.MimeType {
	return f.mimeType
}

func (f *wrapFile) ReadDir(n int) ([]DirEntry, error) {
	entries, err := f.osfile.ReadDir(n)
	if err != nil {
		return nil, err
	}

	domeEntries := make([]DirEntry, 0, len(entries))
	for _, entry := range entries {
		domeEntries = append(domeEntries, WrapDomeEntry(entry))
	}

	return domeEntries, nil
}

type domeEntry struct {
	fs.DirEntry
}

func WrapDomeEntry(e fs.DirEntry) DirEntry {
	return domeEntry{e}
}

// MimeType implements DirEntry
func (e domeEntry) MimeType() entrymodel.MimeType {
	return entrymodel.MimeType(mime.TypeByExtension(path.Ext(e.Name())))
}
