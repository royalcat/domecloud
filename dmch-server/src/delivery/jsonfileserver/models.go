package jsonfileserver

import "dmch-server/src/cfs/media"

type Entry struct {
	Name     string         `json:"name"`
	IsDir    bool           `json:"isDir"`
	MimeType media.MimeType `json:"mimeType"`
}
