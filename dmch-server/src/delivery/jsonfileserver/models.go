package jsonfileserver

import "dmch-server/src/domefs/entrymodel"

type Entry struct {
	Name     string              `json:"name"`
	IsDir    bool                `json:"isDir"`
	MimeType entrymodel.MimeType `json:"mimeType"`
}
