package domefile

import (
	"bytes"
	"dmch-server/src/domefs/entrymodel"
	"io/fs"
	"time"
)

type memoryFile struct {
	name     string
	fileInfo fs.FileInfo
	mimeType entrymodel.MimeType
	*bytes.Reader
}

func NewMemoryFile(name string, mimeType entrymodel.MimeType, data []byte) File {
	return &memoryFile{
		name:     name,
		mimeType: mimeType,
		fileInfo: newMemoryFileInfo(name, int64(len(data))),
		Reader:   bytes.NewReader(data),
	}
}

// Name implements File
func (f *memoryFile) Name() string {
	return f.name
}

// Close implements File
func (*memoryFile) Close() error {
	return nil
}

// MimeType implements File
func (f *memoryFile) MimeType() entrymodel.MimeType {
	return f.mimeType
}

// Stat implements File
func (f *memoryFile) Stat() (fs.FileInfo, error) {
	return f.fileInfo, nil
}

// ReadDir implements File
func (*memoryFile) ReadDir(n int) ([]DirEntry, error) {
	return []DirEntry{}, nil
}

type memoryFileInfo struct {
	name       string
	size       int64
	createTime time.Time
}

func newMemoryFileInfo(name string, size int64) fs.FileInfo {
	return &memoryFileInfo{
		name:       name,
		size:       size,
		createTime: time.Now(),
	}
}

// IsDir implements fs.FileInfo
func (memoryFileInfo) IsDir() bool {
	return false
}

// ModTime implements fs.FileInfo
func (i memoryFileInfo) ModTime() time.Time {
	return i.createTime
}

// Mode implements fs.FileInfo
func (memoryFileInfo) Mode() fs.FileMode {
	return 777
}

// Name implements fs.FileInfo
func (fi memoryFileInfo) Name() string {
	return fi.name
}

// Size implements fs.FileInfo
func (fi memoryFileInfo) Size() int64 {
	return fi.size
}

// Sys implements fs.FileInfo
func (memoryFileInfo) Sys() any {
	return nil
}

func NewMemoryDir(name string, files []DirEntry) File {
	return &memoryDirFile{
		name:     name,
		fileinfo: newMemoryFileInfo(name, 0),
		files:    files,
	}
}

type memoryDirFile struct {
	name     string
	fileinfo fs.FileInfo
	files    []DirEntry
}

// Name implements File
func (f *memoryDirFile) Name() string {
	return f.name
}

// Close implements File
func (*memoryDirFile) Close() error {
	return nil
}

// Read implements File
func (*memoryDirFile) Read([]byte) (int, error) {
	return 0, fs.ErrInvalid
}

// Stat implements File
func (f *memoryDirFile) Stat() (fs.FileInfo, error) {
	return f.fileinfo, nil
}

// Seek implements File
func (*memoryDirFile) Seek(offset int64, whence int) (int64, error) {
	return 0, fs.ErrInvalid
}

// MimeType implements File
func (*memoryDirFile) MimeType() entrymodel.MimeType {
	return entrymodel.MimeTypeDirectory
}

// ReadDir implements File
func (f *memoryDirFile) ReadDir(n int) ([]DirEntry, error) {
	if n <= 0 {
		n = len(f.files)
	}
	return f.files[0:n], nil
}

// type memoryDirEntry struct {
// 	name string
// 	info fs.FileInfo
// }

// // Info implements DirEntry
// func (d *memoryDirEntry) Info() (fs.FileInfo, error) {
// 	return d.info, nil
// }

// // IsDir implements DirEntry
// func (*memoryDirEntry) IsDir() bool {
// 	return false
// }

// // Name implements DirEntry
// func (d *memoryDirEntry) Name() string {
// 	return d.name
// }

// // Type implements DirEntry
// func (*memoryDirEntry) Type() fs.FileMode {
// 	return fs.ModePerm
// }

// // MimeType implements DirEntry
// func (e *memoryDirEntry) MimeType() entrymodel.MimeType {
// 	return entrymodel.MimeType(mime.TypeByExtension(path.Ext(e.name)))
// }

type dirEntryFile struct {
	name string
	File
}

// Name implements DirEntry
func (e *dirEntryFile) Name() string {
	return e.name
}

// Type implements DirEntry
func (e *dirEntryFile) Type() fs.FileMode {
	return fs.ModePerm
}

// IsDir implements DirEntry
func (*dirEntryFile) IsDir() bool {
	return false
}

func (e *dirEntryFile) Info() (fs.FileInfo, error) {
	return e.Info()
}

func WrapFileToDirEntry(name string, file File) DirEntry {
	return &dirEntryFile{
		name: name,
		File: file,
	}
}

type memoryDirInfo struct {
	name       string
	createTime time.Time
}

// IsDir implements fs.FileInfo
func (*memoryDirInfo) IsDir() bool {
	return true
}

// ModTime implements fs.FileInfo
func (i *memoryDirInfo) ModTime() time.Time {
	return i.createTime
}

// Mode implements fs.FileInfo
func (*memoryDirInfo) Mode() fs.FileMode {
	return fs.ModeDir | fs.ModePerm
}

// Name implements fs.FileInfo
func (i *memoryDirInfo) Name() string {
	return i.name
}

// Size implements fs.FileInfo
func (*memoryDirInfo) Size() int64 {
	return 0
}

// Sys implements fs.FileInfo
func (*memoryDirInfo) Sys() any {
	return nil
}

func newMemoryDirInfo(name string) fs.FileInfo {
	return &memoryDirInfo{
		name:       name,
		createTime: time.Now(),
	}
}
