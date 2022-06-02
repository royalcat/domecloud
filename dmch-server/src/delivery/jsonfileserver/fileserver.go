package jsonfileserver

import (
	"dmch-server/src/domefs"
	"dmch-server/src/domefs/domefile"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type fileHandler struct {
	root *domefs.DomeFS
}

func FileServer(root *domefs.DomeFS) http.Handler {
	return &fileHandler{root}
}

func (fh *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	fh.serveEntry(w, r, path.Clean(upath))
}

type condResult int

const (
	condNone condResult = iota
	condTrue
	condFalse
)

func (fh *fileHandler) serveEntry(w http.ResponseWriter, r *http.Request, name string) {
	file, err := fh.root.Open(name)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if stat.IsDir() {
		fh.serveDir(w, r, stat, file)
	} else {
		fh.serveFile(w, r, stat, file)
	}

}

func (fh *fileHandler) serveDir(w http.ResponseWriter, r *http.Request, stat fs.FileInfo, f domefile.File) {
	entries, err := f.ReadDir(0)

	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	jsonEntries := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		jsonEntry := Entry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		}

		jsonEntry.MimeType = entry.MimeType()

		jsonEntries = append(jsonEntries, jsonEntry)
	}

	body, err := json.Marshal(jsonEntries)
	if err != nil {
		logrus.Errorf("Failed to marshal names list: %s", err.Error())
	}
	w.Write(body)
}

func (fh *fileHandler) serveFile(w http.ResponseWriter, r *http.Request, stat fs.FileInfo, f domefile.File) {
	setLastModifiedHeader(w, stat.ModTime())
	done, rangeReq := checkPreconditions(w, r, stat.ModTime())
	if done {
		return
	}

	code := http.StatusOK

	// If Content-Type isn't set, use the file's extension to find it, but
	// if the Content-Type is unset explicitly, do not sniff the type.
	ctypes, haveType := w.Header()["Content-Type"]
	var ctype string
	if !haveType {
		w.Header().Set("Content-Type", string(f.MimeType()))
	} else if len(ctypes) > 0 {
		ctype = ctypes[0]
	}

	size := stat.Size()

	// handle Content-Range header.
	sendSize := size
	var sendContent io.Reader = f
	var content io.ReadSeeker = f
	if size >= 0 {
		ranges, err := parseRange(rangeReq, size)
		if err != nil {
			if err == errNoOverlap {
				w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
			}
			http.Error(w, err.Error(), http.StatusRequestedRangeNotSatisfiable)
			return
		}
		if sumRangesSize(ranges) > size {
			// The total number of bytes in all the ranges
			// is larger than the size of the file by
			// itself, so this is probably an attack, or a
			// dumb client. Ignore the range request.
			ranges = nil
		}
		switch {
		case len(ranges) == 1:
			// RFC 7233, Section 4.1:
			// "If a single part is being transferred, the server
			// generating the 206 response MUST generate a
			// Content-Range header field, describing what range
			// of the selected representation is enclosed, and a
			// payload consisting of the range.
			// ...
			// A server MUST NOT generate a multipart response to
			// a request for a single range, since a client that
			// does not request multiple parts might not support
			// multipart responses."
			ra := ranges[0]
			if _, err := content.Seek(ra.start, io.SeekStart); err != nil {
				http.Error(w, err.Error(), http.StatusRequestedRangeNotSatisfiable)
				return
			}
			sendSize = ra.length
			code = http.StatusPartialContent
			w.Header().Set("Content-Range", ra.contentRange(size))
		case len(ranges) > 1:
			sendSize = rangesMIMESize(ranges, ctype, size)
			code = http.StatusPartialContent

			pr, pw := io.Pipe()
			mw := multipart.NewWriter(pw)
			w.Header().Set("Content-Type", "multipart/byteranges; boundary="+mw.Boundary())
			sendContent = pr
			defer pr.Close() // cause writing goroutine to fail and exit if CopyN doesn't finish.
			go func() {
				for _, ra := range ranges {
					part, err := mw.CreatePart(ra.mimeHeader(ctype, size))
					if err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := content.Seek(ra.start, io.SeekStart); err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := io.CopyN(part, content, ra.length); err != nil {
						pw.CloseWithError(err)
						return
					}
				}
				mw.Close()
				pw.Close()
			}()
		}

		w.Header().Set("Accept-Ranges", "bytes")
		if w.Header().Get("Content-Encoding") == "" {
			w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
		}
	}

	w.WriteHeader(code)

	if r.Method != "HEAD" {
		io.CopyN(w, sendContent, sendSize)
	}
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return fmt.Sprintf("500 Internal Server Error\n%s", err.Error()), http.StatusInternalServerError
}

func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}
