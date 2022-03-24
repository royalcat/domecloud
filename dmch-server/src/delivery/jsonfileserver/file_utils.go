package jsonfileserver

import (
	"encoding/json"
	"net/http"
	"path"
	"sort"

	"github.com/sirupsen/logrus"
)

func (fh *fileHandler) serveDir(w http.ResponseWriter, r *http.Request, name string) {
	entries, err := fh.root.ReadDir(name)

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
		if !jsonEntry.IsDir {
			jsonEntry.MimeType, err = fh.root.MimeType(path.Join(name, jsonEntry.Name))
			if err != nil {
				logrus.Errorf("Error creating mime type: %w", err)
			}
		}

		jsonEntries = append(jsonEntries, jsonEntry)
	}

	body, err := json.Marshal(jsonEntries)
	if err != nil {
		logrus.Errorf("Failed to marshal names list: %s", err.Error())
	}
	w.Write(body)
}
