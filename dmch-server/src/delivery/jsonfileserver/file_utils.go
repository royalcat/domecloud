package jsonfileserver

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"sort"

	"github.com/sirupsen/logrus"
)

type anyDirs interface {
	len() int
	name(i int) string
	isDir(i int) bool
}

type fileInfoDirs []fs.FileInfo

func (d fileInfoDirs) len() int          { return len(d) }
func (d fileInfoDirs) isDir(i int) bool  { return d[i].IsDir() }
func (d fileInfoDirs) name(i int) string { return d[i].Name() }

type dirEntryDirs []fs.DirEntry

func (d dirEntryDirs) len() int          { return len(d) }
func (d dirEntryDirs) isDir(i int) bool  { return d[i].IsDir() }
func (d dirEntryDirs) name(i int) string { return d[i].Name() }

func dirList(w http.ResponseWriter, r *http.Request, f http.File) {
	// Prefer to use ReadDir instead of Readdir,
	// because the former doesn't require calling
	// Stat on every entry of a directory on Unix.
	var dirs anyDirs
	var err error
	if d, ok := f.(fs.ReadDirFile); ok {
		var list dirEntryDirs
		list, err = d.ReadDir(-1)
		dirs = list
	} else {
		var list fileInfoDirs
		list, err = f.Readdir(-1)
		dirs = list
	}

	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs.name(i) < dirs.name(j) })

	serializeDirList(w, dirs)
}

func serializeDirList(w http.ResponseWriter, dirs anyDirs) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	jsonEntries := make([]Entry, 0, dirs.len())
	for i, n := 0, dirs.len(); i < n; i++ {
		entry := Entry{
			Name:  dirs.name(i),
			IsDir: dirs.isDir(i),
		}
		jsonEntries = append(jsonEntries, entry)
	}
	body, err := json.Marshal(jsonEntries)
	if err != nil {
		logrus.Errorf("Failed to marshal names list: %s", err.Error())
	}
	w.Write(body)
}

// func serializeDirList(w http.ResponseWriter, dirs anyDirs) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	fmt.Fprintf(w, "<pre>\n")
// 	for i, n := 0, dirs.len(); i < n; i++ {
// 		name := dirs.name(i)
// 		if dirs.isDir(i) {
// 			name += "/"
// 		}
// 		// name may contain '?' or '#', which must be escaped to remain
// 		// part of the URL path, and not indicate the start of a query
// 		// string or fragment.
// 		url := url.URL{Path: name}
// 		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
// 	}
// 	fmt.Fprintf(w, "</pre>\n")
// }

// var htmlReplacer = strings.NewReplacer(
// 	"&", "&amp;",
// 	"<", "&lt;",
// 	">", "&gt;",
// 	// "&#34;" is shorter than "&quot;".
// 	`"`, "&#34;",
// 	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
// 	"'", "&#39;",
// )