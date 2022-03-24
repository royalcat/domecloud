package jsonfileserver

import "dmch-server/src/domefs/media"

type Entry struct {
	Name     string         `json:"name"`
	IsDir    bool           `json:"isDir"`
	MimeType media.MimeType `json:"mimeType"`
}
