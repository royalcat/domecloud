package jsonfileserver

import "dmch-server/src/domefs/entrymodel"

type Entry struct {
	Name       string              `json:"name"`
	IsListable bool                `json:"isListable"`
	MimeType   entrymodel.MimeType `json:"mimeType"`
}
